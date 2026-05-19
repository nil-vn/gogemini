from datetime import datetime
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy import ForeignKey
from app.utils import db
from .base import BaseModel

class CustomerImage(BaseModel):
    __tablename__ = "customer_image"
    __table_args__ = {"sqlite_autoincrement": True}

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    customer_id: Mapped[int] = mapped_column(ForeignKey("customer.id"), nullable=False)
    file_path: Mapped[str] = mapped_column(nullable=False)
    created_at: Mapped[str] = mapped_column(default=datetime.utcnow)

    customer = relationship("Customer", back_populates="images")

    def __repr__(self):
        return f"<CustomerImage {self.id} for Customer {self.customer_id}>"
