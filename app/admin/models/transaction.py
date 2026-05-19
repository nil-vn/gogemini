from datetime import datetime
from typing import Optional

from sqlalchemy.orm import (
    Mapped,
    mapped_column,
)
from sqlalchemy import Column, Integer, ForeignKey
from sqlalchemy.orm import relationship
from app.utils import db
from .base import BaseModel
from .transaction_car import transaction_car
from sqlalchemy import or_


class Transaction(BaseModel):
    __tablename__ = "transaction"
    __table_args__ = {"sqlite_autoincrement": True}

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    purchase_date: Mapped[Optional[str]]
    selling_price: Mapped[Optional[int]]
    deposit_amount: Mapped[Optional[int]]
    status: Mapped[Optional[str]]
    note: Mapped[Optional[str]]
    created_at: Mapped[str] = mapped_column(default=datetime.utcnow)

    @property
    def status_label(self):
        from app.utils.constants import TransactionStatus
        if not self.status:
            return ""
        try:
            # Handle uppercase enums properly
            return TransactionStatus[self.status.upper()].value
        except (KeyError, AttributeError):
            from flask_babel import lazy_gettext as _
            # Fallback if there's dirty data lying around
            status_map = {
                "paid": getattr(TransactionStatus, "PAID", TransactionStatus.DEPOSITED).value if hasattr(TransactionStatus, "PAID") else _("Paid"),
                "wait2pay": getattr(TransactionStatus, "WAIT_TO_PAY", TransactionStatus.DEPOSITED).value if hasattr(TransactionStatus, "WAIT_TO_PAY") else _("Wait to pay"),
                "deposited": getattr(TransactionStatus, "DEPOSITED", TransactionStatus.DEPOSITED).value if hasattr(TransactionStatus, "DEPOSITED") else _("Deposited")
            }
            return status_map.get(self.status.lower(), self.status)

    customer_id: Mapped[Optional[int]] = mapped_column(
        ForeignKey("customer.id", ondelete="CASCADE"), nullable=True
    )

    # Quan hệ với Customer
    customer = relationship("Customer", back_populates="transactions")

    # Quan hệ nhiều-nhiều với Car
    cars = relationship("Car", secondary=transaction_car, back_populates="transactions")

    # Quan hệ 1-Nhiều với TransactionItem
    items = relationship("TransactionItem", back_populates="transaction", cascade="all, delete-orphan")

    @property
    def total_other_amount(self):
        return sum(item.price for item in self.items if item.price)

    @property
    def total_amount(self):
        return (self.selling_price or 0) + self.total_other_amount

    def __repr__(self):
        return f"<Transaction {self.name}>"

    @classmethod
    def field_names(cls):
        """Return all column names for this model."""
        return [c.name for c in cls.__table__.columns]

    @classmethod
    def from_form(cls, form):
        """Create instance from a form, filtering only valid fields."""
        data = {k: v for k, v in form.data.items() if k in cls.field_names()}
        return cls(**data)

    @classmethod
    def get_all(cls):
        return cls.query.all()

    @classmethod
    def get_by_id(cls, qid):
        return cls.query.get(qid)

    @classmethod
    def search(cls, query):
        return cls.query.filter(
            or_(cls.status.ilike(f"%{query}%"), cls.note.ilike(f"%{query}%"))
        ).all()

    @classmethod
    def get_by_statuses(cls, statuses: list):
        return cls.query.filter(cls.status.in_(statuses)).all()

