from datetime import datetime
from types import SimpleNamespace

from app.admin.models import Customer, Transaction, Car


def get_all_customers():
    customers = Customer.get_all()
    return customers

def get_recently_count(model):
    now = datetime.utcnow()
    # tháng hiện tại
    start_current_month = datetime(now.year, now.month, 1)
    # tháng trước
    if now.month == 1:
        start_prev_month = datetime(now.year - 1, 12, 1)
    else:
        start_prev_month = datetime(now.year, now.month - 1, 1)

    # số user tháng hiện tại
    current_count = model.current_count(start_current_month)

    # số user tháng trước
    prev_count = model.prev_count(start_current_month, start_prev_month)

    # tính % tăng trưởng
    growth_rate = 0
    if prev_count > 0:
        growth_rate = ((current_count - prev_count) / prev_count) * 100

    return SimpleNamespace(
        total=current_count,
        monthly_increase=round(growth_rate, 2)
    )


def get_revenue_last_6_months():
    now = datetime.utcnow()
    results = []
    # from current month back to 5 months ago
    for i in range(5, -1, -1):
        # calculate year and month
        y = now.year
        m = now.month - i
        while m <= 0:
            m += 12
            y -= 1
            
        start_date = datetime(y, m, 1)
        if m == 12:
            end_date = datetime(y + 1, 1, 1)
        else:
            end_date = datetime(y, m + 1, 1)
            
        # Get transactions in this month
        transactions = Transaction.query.filter(
            Transaction.created_at >= start_date,
            Transaction.created_at < end_date
        ).all()
        
        month_revenue = sum(t.total_amount for t in transactions if getattr(t, 'total_amount', None))
        results.append({
            "month": f"{m:02d}/{y}",
            "revenue": month_revenue
        })
    return results

def get_car_status_stats():
    cars = Car.get_all()
    stats = {}
    for car in cars:
        status = car.car_situation or 'Unknown'
        if hasattr(status, 'value'):
            status = status.value
            
        if status not in stats:
            stats[status] = 0
        stats[status] += 1
    
    # Format for chart
    labels = list(stats.keys())
    series = list(stats.values())
    return {
        "labels": labels,
        "series": series
    }

def get_metrics():
    return SimpleNamespace(
        customers=get_recently_count(Customer),
        transactions=get_recently_count(Transaction),
        cars=get_recently_count(Car),
        revenue_6_months=get_revenue_last_6_months(),
        car_stats=get_car_status_stats()
    )