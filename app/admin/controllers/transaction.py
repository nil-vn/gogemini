from flask_babel import _
from flask import request, redirect
from flask import render_template, flash, url_for

from app.admin.models import Car, Customer, Transaction
from app.admin.services.create_or_updating import create_transaction_from_form
from app.admin.services.forms import (
    TransactionForm,
)
from app.utils.db import db
from flask_login import login_required
from . import routes
from ...utils.constants import TransactionStatus


@routes.route("/transactions")
@login_required
def transactions():
    all_transactions = Transaction.get_all()
    deposited_transactions = Transaction.get_by_statuses([TransactionStatus.DEPOSITED.name])
    paid_transactions = Transaction.get_by_statuses([TransactionStatus.PAID.name])
    
    total_revenue = sum(t.total_amount or 0 for t in all_transactions if t.status in [TransactionStatus.PAID.name, TransactionStatus.DEPOSITED.name])
    paid_revenue = sum(t.total_amount or 0 for t in paid_transactions)
    deposited_amount = sum(t.deposit_amount or 0 for t in deposited_transactions)

    return render_template(
        "transactions.html",
        transactions=all_transactions,
        deposited_transactions=deposited_transactions,
        paid_transactions=paid_transactions,
        total_revenue=total_revenue,
        paid_revenue=paid_revenue,
        deposited_amount=deposited_amount
    )


@routes.route("/transaction/new", methods=["GET", "POST"])
@login_required
def transaction_new():
    form = TransactionForm()
    customer_id = request.args.get("customer_id")
    car_id = request.args.get("car_id")

    # Prefill nếu có từ query param
    if customer_id:
        form.customer_id.data = customer_id
    if car_id:
        form.car_id.data = car_id

    if request.method == "POST":
        if form.validate_on_submit():
            try:
                tx = create_transaction_from_form(form)
                from app.admin.models.transaction_item import TransactionItem
                item_names = request.form.getlist('item_name[]')
                item_prices = request.form.getlist('item_price[]')
                for name, price in zip(item_names, item_prices):
                    if name and str(price).strip():
                        item = TransactionItem(transaction_id=tx.id, name=name, price=int(str(price).replace(',','')))
                        db.session.add(item)
                db.session.commit()
                flash(_("Transaction created successfully!"), "success")
                return redirect(url_for("admin_routes.transactions"))
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error creating transaction: {e}"), "danger")
            finally:
                return redirect(url_for('admin_routes.transaction_new'))
        elif form.errors:
            for err_code, err_content in form.errors.items():
                for e in err_content:
                    flash(_(f"{err_code}: {e}"), "danger")
            return redirect(url_for('admin_routes.transaction_new'))

    cars = Car.get_all()
    customers = Customer.get_all()
    transaction_status = TransactionStatus
    return render_template(
        "transaction_new.html", cars=cars, customers=customers, form=form, transaction_status=transaction_status
    )


@routes.route("/transaction/<transaction_id>", methods=["GET", "POST"])
@login_required
def transaction_detail(transaction_id):
    form = TransactionForm()
    transaction = Transaction.get_by_id(transaction_id)
    if not transaction:
        flash(_("Transaction not found."), "danger")
        return render_template("404.html"), 404

    if request.method == "POST":
        form = TransactionForm(obj=transaction)
        if form.validate_on_submit():
            try:
                form.populate_obj(transaction)
                # clear old items
                for item in transaction.items:
                    db.session.delete(item)
                # add new items
                from app.admin.models.transaction_item import TransactionItem
                item_names = request.form.getlist('item_name[]')
                item_prices = request.form.getlist('item_price[]')
                for name, price in zip(item_names, item_prices):
                    if name and str(price).strip():
                        item_price = int(str(price).replace(',', ''))
                        item = TransactionItem(transaction_id=transaction.id, name=name, price=item_price)
                        db.session.add(item)
                        
                db.session.commit()
                flash(_("Transaction updated successfully!"), "success")
            except Exception as e:
                db.session.rollback()
                flash(_(f"Error updating transaction: {e}"), "danger")
    customers = Customer.get_all()
    cars = Car.get_all()
    transaction_status = TransactionStatus
    return render_template(
        "transaction_detail.html",
        transaction=transaction,
        form=form,
        customers=customers,
        cars=cars,
        transaction_status=transaction_status
    )


@routes.route("/transaction/<int:transaction_id>/delete", methods=["GET"])
@login_required
def delete_transaction(transaction_id):
    transaction = Transaction.get_by_id(transaction_id)
    if not transaction:
        flash(_("Transaction not found."), "danger")
        return redirect(url_for("admin_routes.transactions"))

    try:
        db.session.delete(transaction)
        db.session.commit()
        flash(_("Transaction deleted successfully!"), "success")
    except Exception as e:
        db.session.rollback()  # rollback nếu lỗi
        flash(_(f"Error deleting transaction: {e}"), "danger")

    return redirect(url_for("admin_routes.transactions"))
