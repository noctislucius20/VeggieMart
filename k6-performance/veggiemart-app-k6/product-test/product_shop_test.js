import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
    vus: 25,
    duration: "1m",
};

export default function () {
    const BASE_URL = "http://localhost:8081";
    const query = {
        page: 1,
        limit: 10,
        order_by: "created_at",
        order_type: "desc",
        start_price: 0,
        end_price: 0,
        search: "",
    };
    const res = http.get(
        `${BASE_URL}/products/shop?` +
            `page=${query.page}&limit=${query.limit}` +
            `&order_by=${query.order_by}&order_type=${query.order_type}` +
            `&start_price=${query.start_price}&end_price=${query.end_price}` +
            `&search=${query.search}`,
    );

    check(res, {
        "status is 200": (r) => r.status === 200,
        "message is success": (r) => JSON.parse(r.body).message === "success",
        "has pagination": (r) => JSON.parse(r.body).pagination !== null,
        "response is JSON": (r) =>
            r.headers["Content-Type"] === "application/json",
        "data is array": (r) => Array.isArray(JSON.parse(r.body).data),
    });

    sleep(1);
}
