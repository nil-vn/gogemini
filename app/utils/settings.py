from app.utils.db import db
from app.admin.models.base import Configuration


CURRENCY_SYMBOLS = {
    "VND": "đ",
    "JPY": "¥",
}

DEFAULTS = {
    "currency": "JPY",
    "theme": "dark",
    "language": "vi",
}


def get_setting(key: str) -> str:
    """Get a system setting value by key, returning the default if not found."""
    row = Configuration.query.filter_by(key=key).first()
    if row:
        return row.value
    return DEFAULTS.get(key, "")


def set_setting(key: str, value: str) -> None:
    """Set a system setting value by key (insert or update)."""
    row = Configuration.query.filter_by(key=key).first()
    if row:
        row.value = value
    else:
        row = Configuration(key=key, value=value)
        db.session.add(row)
    db.session.commit()


def get_currency_symbol() -> str:
    """Return the currency symbol based on the current currency setting."""
    currency = get_setting("currency")
    return CURRENCY_SYMBOLS.get(currency, currency)


def format_currency(value) -> str:
    """Jinja filter: format a number with the current currency symbol.
    Example: 1500000 -> '1,500,000 đ' or '1,500,000 ¥'
    """
    if value is None:
        return "-"
    try:
        formatted = "{:,.0f}".format(float(value))
    except (ValueError, TypeError):
        return str(value)
    return f"{formatted} {get_currency_symbol()}"
