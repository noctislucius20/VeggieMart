import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
    vus: 30,
    duration: "30s",
};

export default function () {
    const lat = -6.175393;
    const lng = 106.827153;

    const url = `http://localhost:8082/auth/orders?lat=${lat}&lng=${lng}`;

    const payload = JSON.stringify({
        buyer_id: 3,
        order_date: "2026-04-06",
        total_amount: 10000,
        shipping_type: "Delivery",
        remarks: "",
        order_time: "13:00:00",
        order_details: [
            {
                product_id: 5,
                quantity: 1,
            },
        ],
    });

    const params = {
        headers: {
            "Content-Type": "application/json",
            Authorization:
                "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzU2MzA4OTUsImlzcyI6InNlY3JldHp6enp6IiwidXNlcl9pZCI6M30.A_W8J694nlv42iAUQdQ2jYJTZVvVnUcr-NBNThGjQeo",
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        "status is 200": (r) => r.status === 201,
        "response says success": (r) => r.json("message") === "success",
    });

    sleep(1);
}
