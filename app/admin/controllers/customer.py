from flask_babel import _
from flask import request, redirect, jsonify
from flask import render_template, flash, url_for
import os
import uuid
from werkzeug.utils import secure_filename
from flask import current_app

from app.admin.models import Car, Customer, CustomerImage
from app.admin.services.create_or_updating import create_transaction_from_form
from app.admin.services.forms import (
    CustomerForm,
    TransactionForm,
)
from app.utils.db import db
from flask_login import login_required
from . import routes


@routes.route("/customers")
@login_required
def customers():
    all_customers = Customer.get_all()
    active_customers = [c for c in all_customers if c.transactions]
    return render_template("customers.html", customers=all_customers, active_customers=active_customers)


@routes.route("/customer/new", methods=["GET", "POST"])
@login_required
def customer_new():
    form = CustomerForm()
    if request.method == "POST":
        if form.validate_on_submit():
            # Add to session and commit
            new_customer = Customer.from_form(form)
            try:
                db.session.add(new_customer)
                db.session.flush() # get new customer ID
                
                # handle image uploads
                if 'images' in request.files:
                    files = request.files.getlist('images')
                    upload_folder = os.path.join(current_app.root_path, '..', 'static', 'uploads', 'customers')
                    os.makedirs(upload_folder, exist_ok=True)
                    for file in files:
                        if file and file.filename:
                            original_filename = secure_filename(file.filename)
                            save_name = f"{uuid.uuid4().hex}_{original_filename}"
                            file_path = os.path.join(upload_folder, save_name)
                            file.save(file_path)
                            
                            customer_image = CustomerImage(customer_id=new_customer.id, file_path=f"uploads/customers/{save_name}")
                            db.session.add(customer_image)

                db.session.commit()
                flash(_("Customer added successfully!"), "success")
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error adding car: {e}"), "danger")
            finally:
                return redirect(url_for('admin_routes.customer_new'))
        elif form.errors:
            for err_code, err_content in form.errors.items():
                for e in err_content:
                    flash(_(f"{err_code}: {e}"), "danger")
            return redirect(url_for('admin_routes.customer_new'))
    return render_template("customer_new.html", form=form)


@routes.route("/customer/<int:customer_id>", methods=["GET", "POST"])
@login_required
def customer_detail(customer_id):
    form = CustomerForm()
    customer = Customer.get_by_id(customer_id)
    if not customer:
        flash(_("Customer not found."), "danger")
        return render_template("404.html"), 404

    if request.method == "POST":
        form = CustomerForm(obj=customer)
        if form.validate_on_submit():
            try:
                form.populate_obj(customer)
                
                # handle image uploads
                if 'images' in request.files:
                    files = request.files.getlist('images')
                    upload_folder = os.path.join(current_app.root_path, '..', 'static', 'uploads', 'customers')
                    os.makedirs(upload_folder, exist_ok=True)
                    for file in files:
                        if file and file.filename:
                            original_filename = secure_filename(file.filename)
                            save_name = f"{uuid.uuid4().hex}_{original_filename}"
                            file_path = os.path.join(upload_folder, save_name)
                            file.save(file_path)
                            
                            customer_image = CustomerImage(customer_id=customer.id, file_path=f"uploads/customers/{save_name}")
                            db.session.add(customer_image)

                db.session.commit()
                flash(_("Customer updated successfully!"), "success")
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error updating Customer: {e}"), "danger")
    cars = Car.get_all()
    return render_template(
        "customer_detail.html", customer=customer, form=form, cars=cars
    )


@routes.route("/customer/<int:customer_id>/purchase/new", methods=["GET", "POST"])
@login_required
def add_customer_purchase(customer_id):
    form = TransactionForm()
    form.customer_id.data = int(customer_id)
    # Prefill customer_id
    if request.method == "POST":
        if form.validate_on_submit():
            try:
                create_transaction_from_form(form)
                flash(_("Purchase added for customer successfully!"), "success")
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error creating transaction: {e}"), "danger")
        elif form.errors:
            for err_code, err_content in form.errors.items():
                for e in err_content:
                    flash(_(f"{err_code}: {e}"), "danger")
    return redirect(url_for("admin_routes.customer_detail", customer_id=customer_id))


@routes.route("/customer/<int:customer_id>/delete", methods=["GET"])
@login_required
def delete_customer(customer_id):
    customer = Customer.get_by_id(customer_id)
    if not customer:
        flash(_("Customer not found."), "danger")
        return redirect(url_for("admin_routes.customers"))

    try:
        db.session.delete(customer)
        db.session.commit()
        flash(_("Customer deleted successfully!"), "success")
    except Exception as e:
        db.session.rollback()  # rollback nếu lỗi
        flash(_(f"Error deleting car: {e}"), "danger")

    return redirect(url_for("admin_routes.customers"))

@routes.route("/customer/image/<int:image_id>/delete", methods=["POST"])
@login_required
def delete_customer_image(image_id):
    image = CustomerImage.query.get(image_id)
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
