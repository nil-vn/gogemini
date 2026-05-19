from flask_babel import _
from flask import request, redirect, jsonify
from flask import render_template, flash, url_for

from app.admin.models import Car, Customer, CarImage
import os
import uuid
from werkzeug.utils import secure_filename
from flask import current_app
from app.admin.services.forms import (
    CarForm,
    TransactionForm,
)
from app.utils.db import db
from flask_login import login_required
from . import routes
from ..services.create_or_updating import create_transaction_from_form
from ...utils.constants import CarStatus, CarBranches, CarSituation, TransactionStatus


@routes.route("/cars")
@login_required
def cars():
    all_cars = Car.get_all()
    # 2. Có sẵn (có tính cả chờ giao)
    available_cars = Car.get_by_statuses([CarStatus.AVAILABLE.name, CarStatus.AWAITING_DELIVERY.name])
    # 3. Đang chờ giao
    awaiting_delivery_cars = Car.get_by_statuses([CarStatus.AWAITING_DELIVERY.name])
    # 4. Đã bán
    sold_cars = Car.get_by_statuses([CarStatus.SOLD.name])
    
    # 5. Sơn sửa rồi
    refurbished_cars = Car.get_by_situations([CarSituation.REFURBISHED.name])
    # 6. Chưa sơn sửa
    not_refurbished_cars = Car.get_by_situations([CarSituation.NOT_REFURBISHED.name])
    # 7. Sơn sửa rồi - chưa vệ sinh
    refurbished_pending_cleaning_cars = Car.get_by_situations([CarSituation.REFURBISHED_PENDING_CLEANING.name])

    return render_template(
        "cars.html",
        cars=all_cars,
        available_cars=available_cars,
        awaiting_delivery_cars=awaiting_delivery_cars,
        sold_cars=sold_cars,
        refurbished_cars=refurbished_cars,
        not_refurbished_cars=not_refurbished_cars,
        refurbished_pending_cleaning_cars=refurbished_pending_cleaning_cars,
    )


@routes.route("/car/new", methods=["GET", "POST"])
@login_required
def car_new():
    form = CarForm()
    car_status = CarStatus
    car_branches = CarBranches
    car_situation = CarSituation
    if request.method == "POST":
        if form.validate_on_submit():
            # Add to session and commit
            new_car = Car.from_form(form)
            try:
                db.session.add(new_car)
                db.session.flush() # get new car ID
                
                # handle image uploads
                if 'images' in request.files:
                    files = request.files.getlist('images')
                    upload_folder = os.path.join(current_app.root_path, '..', 'static', 'uploads', 'cars')
                    os.makedirs(upload_folder, exist_ok=True)
                    for file in files:
                        if file and file.filename:
                            original_filename = secure_filename(file.filename)
                            save_name = f"{uuid.uuid4().hex}_{original_filename}"
                            file_path = os.path.join(upload_folder, save_name)
                            file.save(file_path)
                            
                            car_image = CarImage(car_id=new_car.id, file_path=f"uploads/cars/{save_name}")
                            db.session.add(car_image)

                db.session.commit()
                # Nếu đến đây, commit đã thành công
                flash(_("Car added successfully!"), "success")
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error adding car: {e}"), "danger")
            finally:
                return redirect(url_for('admin_routes.car_new'))
        elif form.errors:
            for err_code, err_content in form.errors.items():
                for e in err_content:
                    flash(_(f"{err_code}: {e}"), "danger")
            return redirect(url_for('admin_routes.car_new'))
    return render_template(
        "car_new.html",
        form=form,
        car_status=car_status,
        car_branches=car_branches,
        car_situation=car_situation
    )


@routes.route("/car/<car_id>", methods=["GET", "POST"])
@login_required
def car_detail(car_id):
    form = CarForm()
    car = Car.get_by_id(car_id)
    if not car:
        flash(_("Car not found."), "danger")
        return render_template("404.html"), 404

    if request.method == "POST":
        form = CarForm(obj=car)
        if form.validate_on_submit():
            try:
                form.populate_obj(car)
                
                # handle image uploads
                if 'images' in request.files:
                    files = request.files.getlist('images')
                    upload_folder = os.path.join(current_app.root_path, '..', 'static', 'uploads', 'cars')
                    os.makedirs(upload_folder, exist_ok=True)
                    for file in files:
                        if file and file.filename:
                            original_filename = secure_filename(file.filename)
                            save_name = f"{uuid.uuid4().hex}_{original_filename}"
                            file_path = os.path.join(upload_folder, save_name)
                            file.save(file_path)
                            
                            car_image = CarImage(car_id=car.id, file_path=f"uploads/cars/{save_name}")
                            db.session.add(car_image)

                db.session.commit()
                flash(_("Car updated successfully!"), "success")
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error updating car: {e}"), "danger")
            finally:
                return redirect(url_for('admin_routes.car_detail', car_id=car_id))
        elif form.errors:
            for err_code, err_content in form.errors.items():
                for e in err_content:
                    flash(_(f"{err_code}: {e}"), "danger")
            return redirect(url_for('admin_routes.car_detail', car_id=car_id))
    customers = Customer.get_all()
    car_status = CarStatus
    car_branches = CarBranches
    car_situation = CarSituation
    transaction_status = TransactionStatus
    return render_template(
        "car_detail.html",
        car=car,
        form=form,
        customers=customers,
        car_status=car_status,
        car_branches=car_branches,
        car_situation=car_situation,
        transaction_status=transaction_status
    )


@routes.route("/car/<int:car_id>/delete", methods=["GET"])
@login_required
def delete_car(car_id):
    car = Car.get_by_id(car_id)
    if not car:
        flash(_("Car not found."), "danger")
        return redirect(url_for("admin_routes.cars"))

    try:
        db.session.delete(car)
        db.session.commit()
        flash(_("Car deleted successfully!"), "success")
    except Exception as e:
        db.session.rollback()
        flash(_(f"Error deleting car: {e}"), "danger")

    return redirect(url_for("admin_routes.cars"))


@routes.route("/car/<int:car_id>/purchase/new", methods=["GET", "POST"])
@login_required
def add_car_purchase(car_id):
    form = TransactionForm()
    customers = Customer.get_all()
    car = Car.get_by_id(car_id)
    form.car_id.data = int(car_id)  # Prefill car_id
    if request.method == "POST":
        if form.validate_on_submit():
            try:
                create_transaction_from_form(form)
                flash(_("Purchase added for car successfully!"), "success")
            except Exception as e:
                flash(_(f"Error creating transaction: {e}"), "danger")
        elif form.errors:
            for err_code, err_content in form.errors.items():
                for e in err_content:
                    flash(_(f"{err_code}: {e}"), "danger")
    return redirect(url_for("admin_routes.car_detail", car_id=car.id))


@routes.route("/car/image/<int:image_id>/delete", methods=["POST"])
@login_required
def delete_car_image(image_id):
    image = CarImage.query.get(image_id)
    if not image:
        return jsonify({"success": False, "error": "Image not found"}), 404
        
    try:
        # Delete file from filesystem
        file_path = os.path.join(current_app.root_path, '..', 'static', image.file_path)
        if os.path.exists(file_path):
            os.remove(file_path)
            
        # Delete from database
        db.session.delete(image)
        db.session.commit()
        return jsonify({"success": True})
    except Exception as e:
        db.session.rollback()
        return jsonify({"success": False, "error": str(e)}), 500
