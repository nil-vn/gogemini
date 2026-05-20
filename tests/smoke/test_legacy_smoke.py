import io
import os
from datetime import datetime

import pytest

from app import create_app
from app.utils import db
from app.admin.models import User, Car, Customer, Transaction


@pytest.fixture()
def client():
    os.makedirs("instance", exist_ok=True)
    app = create_app("config.dev.DevConfig")
    app.config.update(TESTING=True, WTF_CSRF_ENABLED=False)

    with app.app_context():
        db.drop_all()
        db.create_all()

        admin = User(
            username="smoke_admin",
            email="smoke_admin@example.com",
            role="admin",
            status="Active",
        )
        admin.set_password("smoke-pass-123")
        db.session.add(admin)
        db.session.commit()

    with app.test_client() as test_client:
        yield test_client


def test_legacy_smoke_baseline(client):
    # login flow
    resp = client.post(
        "/admin/login",
        data={"username": "smoke_admin", "password": "smoke-pass-123", "remember": "y"},
        follow_redirects=True,
    )
    assert resp.status_code == 200

    # users CRUD smoke: create + list + update + delete
    create_user = client.post(
        "/admin/user/new",
        data={
            "username": "smoke_user",
            "email": "smoke_user@example.com",
            "password": "smoke-pass-456",
            "confirm_password": "smoke-pass-456",
            "role": "staff",
            "status": "Active",
        },
        follow_redirects=True,
    )
    assert create_user.status_code == 200
    users_page = client.get("/admin/users")
    assert users_page.status_code == 200

    with create_app("config.dev.DevConfig").app_context():
        user = User.query.filter_by(username="smoke_user").first()
        assert user is not None
        uid = user.id

    update_user = client.post(
        f"/admin/user/{uid}",
        data={
            "username": "smoke_user_updated",
            "email": "smoke_user_updated@example.com",
            "password": "",
            "role": "staff",
            "status": "Active",
        },
        follow_redirects=True,
    )
    assert update_user.status_code == 200
    delete_user = client.get(f"/admin/user/{uid}/delete", follow_redirects=True)
    assert delete_user.status_code == 200

    # cars CRUD smoke + upload
    car_create = client.post(
        "/admin/car/new",
        data={
            "name": "Smoke Car",
            "vin": "VIN-SMOKE-001",
            "model": "Model S",
            "branch": "Smoke",
            "color": "Black",
            "traded_company": "SmokeCo",
            "imported_date": "2026-01-01",
            "inspection_from": "2026-01-01",
            "inspection_to": "2026-12-31",
            "year_of_manufacture": "2026",
            "purchase_price": "1000",
            "selling_price": "1500",
            "status": "Available",
            "note": "smoke",
            "license_plate_no": "SMK-001",
            "car_situation": "Good",
            "gallery": (io.BytesIO(b"fake image bytes"), "smoke-car.jpg"),
        },
        content_type="multipart/form-data",
        follow_redirects=True,
    )
    assert car_create.status_code == 200

    cars_page = client.get("/admin/cars")
    assert cars_page.status_code == 200

    # customers CRUD smoke
    customer_create = client.post(
        "/admin/customer/new",
        data={
            "name": "Smoke Customer",
            "gender": "male",
            "birth_day": "1990-01-01",
            "facebook": "fb/smoke",
            "phone": "0123456",
            "address": "Smoke Street",
            "lead_source": "web",
            "status": "active",
            "note": "smoke",
        },
        follow_redirects=True,
    )
    assert customer_create.status_code == 200
    customers_page = client.get("/admin/customers")
    assert customers_page.status_code == 200

    # transactions CRUD smoke (create + list)
    with create_app("config.dev.DevConfig").app_context():
        car = Car.query.filter_by(vin="VIN-SMOKE-001").first()
        customer = Customer.query.filter_by(name="Smoke Customer").first()
        assert car and customer

    tx_create = client.post(
        "/admin/transaction/new",
        data={
            "customer_id": str(customer.id),
            "car_id": str(car.id),
            "purchase_date": datetime.utcnow().strftime("%Y-%m-%d"),
            "selling_price": "2000",
            "deposit_amount": "200",
            "status": "pending",
            "note": "smoke",
        },
        follow_redirects=True,
    )
    assert tx_create.status_code == 200
    tx_page = client.get("/admin/transactions")
    assert tx_page.status_code == 200

    # search flow
    search_resp = client.get("/admin/search?q=Smoke")
    assert search_resp.status_code == 200
    assert b"Smoke" in search_resp.data
