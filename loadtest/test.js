import http from "k6/http";

export default function () {
  const url = "http://localhost:8080/journal";

  const payload = JSON.stringify({
    idempotencyKey: "183",
    from: "jesse",
    to: "jan",
    amount: {
      value: Math.floor(Math.random() * 100) + 1,
      currency: "EUR",
    },
    metadata: {
      orderReference: "39123912",
    },
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  http.post(url, payload, params);
}
