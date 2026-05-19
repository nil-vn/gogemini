from datetime import datetime
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy import ForeignKey
from app.utils import db
from .base import BaseModel

class CarImage(BaseModel):
    __tablename__ = "car_image"
    __table_args__ = {"sqlite_autoincrement": True}

    id: Mapped[int] = mapped_column(primary_key=True, autoincrement=True)
    car_id: Mapped[int] = mapped_column(ForeignKey("car.id"), nullable=False)
    file_path: Mapped[str] = mapped_column(nullable=False)
    created_at: Mapped[str] = mapped_column(default=datetime.utcnow)

    car = relationship("Car", back_populates="images")

    def __repr__(self):
        return f"<CarImage {self.id} for Car {self.car_id}>"
