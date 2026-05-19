from types import SimpleNamespace

from flask import redirect, request
from flask import flash, url_for
from flask_babel import _ as _t

from app.admin.models import User, Customer
from app.admin.services.forms import (
    RegisterForm,
)
from app.utils.db import db
from werkzeug.security import generate_password_hash

from flask import render_template
from flask_login import login_required
from . import routes
from ..services.analytics import get_metrics
from app.utils.settings import get_setting, set_setting


@routes.route("/")
@routes.route("/dashboard")
# @login_required
def dashboard():
    metric = get_metrics()
    return render_template("dashboard.html", metric=metric)

@routes.route("/system", methods=["GET", "POST"])
# @login_required
def system():
    if request.method == "POST":
        set_setting("currency", request.form.get("currency", "JPY"))
        set_setting("theme", request.form.get("theme", "dark"))
        set_setting("language", request.form.get("language", "vi"))
        flash(_t("Settings saved successfully!"), "success")
        return redirect(url_for("admin_routes.system"))

    settings = {
        "currency": get_setting("currency"),
        "theme": get_setting("theme"),
        "language": get_setting("language"),
    }
    return render_template("system.html", settings=settings)
