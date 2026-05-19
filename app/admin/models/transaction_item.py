from typing import Optional
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy import ForeignKey
from app.utils.db import db
from .base import BaseModel

class TransactionItem(BaseModel):
    __tablename__ = "transaction_item"
    __table_args__ = {"sqlite_autoincrement": True}

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    transaction_id: Mapped[int] = mapped_column(ForeignKey("transaction.id", ondelete="CASCADE"))
    
    name: Mapped[str]
    price: Mapped[Optional[int]]

    transaction = relationship("Transaction", back_populates="items")

    def __repr__(self):
        return f"<TransactionItem {self.name}>"

    @classmethod
    def from_form(cls, form):
        return cls()
