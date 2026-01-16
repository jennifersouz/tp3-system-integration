import os
import time
import csv
import io
import random
from datetime import datetime

import boto3
from botocore.config import Config


def env(name: str) -> str:
    v = os.getenv(name)
    if not v:
        raise RuntimeError(f"Missing env var: {name}")
    return v


S3_ENDPOINT = env("S3_ENDPOINT")
S3_REGION = env("S3_REGION")
S3_ACCESS_KEY_ID = env("S3_ACCESS_KEY_ID")
S3_SECRET_ACCESS_KEY = env("S3_SECRET_ACCESS_KEY")
S3_BUCKET = env("S3_BUCKET")

PERIOD_SECONDS = int(os.getenv("CRAWLER_PERIOD_SECONDS", "30"))


def make_s3_client():
    # Supabase S3 normalmente exige path-style
    return boto3.client(
        "s3",
        endpoint_url=S3_ENDPOINT,
        region_name=S3_REGION,
        aws_access_key_id=S3_ACCESS_KEY_ID,
        aws_secret_access_key=S3_SECRET_ACCESS_KEY,
        config=Config(s3={"addressing_style": "path"}),
    )


def generate_csv_bytes() -> bytes:
    # Gera um CSV compatível com o teu schema (Superstore-like)
    header = [
        "Row ID","Order ID","Order Date","Ship Date","Ship Mode","Customer ID","Customer Name","Segment",
        "Country","City","State","Postal Code","Region","Retail Sales People","Product ID","Category",
        "Sub-Category","Product Name","Returned","Sales","Quantity","Discount","Profit"
    ]

    ship_modes = ["Second Class", "Standard Class", "First Class"]
    segments = ["Consumer", "Corporate", "Home Office"]
    states = ["Kentucky", "California", "Florida", "New York"]
    regions = {"Kentucky":"South", "California":"West", "Florida":"South", "New York":"East"}
    categories = [("Furniture","Chairs"), ("Office Supplies","Storage"), ("Technology","Phones")]

    buf = io.StringIO()
    w = csv.writer(buf)

    w.writerow(header)

    now = datetime.utcnow()
    order_id = f"TP3-{now.strftime('%Y%m%d%H%M%S')}"
    order_date = now.strftime("%Y-%m-%d")
    ship_date = now.strftime("%Y-%m-%d")

    for i in range(1, 6):  # 5 linhas por ficheiro
        state = random.choice(states)
        category, sub = random.choice(categories)

        sales = round(random.uniform(10, 1200), 4)
        qty = random.randint(1, 8)
        discount = round(random.choice([0.0, 0.1, 0.2, 0.3]), 2)
        profit = round(sales * (random.uniform(-0.3, 0.3)), 4)

        row = [
            i,
            order_id,
            order_date,
            ship_date,
            random.choice(ship_modes),
            f"CUST-{random.randint(10000,99999)}",
            f"Cliente {random.randint(1,200)}",
            random.choice(segments),
            "United States",
            "CityX",
            state,
            str(random.randint(10000, 99999)),
            regions[state],
            "SalesPerson X",
            f"PROD-{random.randint(1000,9999)}",
            category,
            sub,
            f"Produto {random.randint(1,999)}",
            "Not",
            sales,
            qty,
            discount,
            profit
        ]
        w.writerow(row)

    return buf.getvalue().encode("utf-8")


def main():
    s3 = make_s3_client()

    while True:
        content = generate_csv_bytes()
        key = f"input/orders_{datetime.utcnow().strftime('%Y%m%d_%H%M%S')}.csv"

        s3.put_object(
            Bucket=S3_BUCKET,
            Key=key,
            Body=content,
            ContentType="text/csv"
        )

        print(f"[crawler] uploaded: s3://{S3_BUCKET}/{key} ({len(content)} bytes)")
        time.sleep(PERIOD_SECONDS)


if __name__ == "__main__":
    main()
