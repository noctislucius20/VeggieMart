import http from "k6/http";
import { check, sleep } from "k6";

export let options = {
    vus: 30,
    duration: "30s",
};

export default function () {
    const url = "http://localhost:8080/users/signin";

    const payload = JSON.stringify({
        email: "test@mail.com",
        password: "scipio123",
    });

    const params = {
        headers: {
            "Content-Type": "application/json",
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        "status is 200": (r) => r.status === 200,
        "response says success": (r) => r.json("message") === "success",
    });

    sleep(1);
}
