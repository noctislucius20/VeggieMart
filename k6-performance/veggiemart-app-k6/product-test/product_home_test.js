import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    vus: 20,
    duration: "1m",
};

export default function () {
    const BASE_URL = "http://localhost:8081";
    const res = http.get(`${BASE_URL}/products/home`);

    check(res, {
        "status is 200": (r) => r.status === 200,
        "response is JSON": (r) =>
            r.headers["Content-Type"] === "application/json",
        "has categories data": (r) => JSON.parse(r.body).data.length >= 0,
    });

    sleep(1);
}
