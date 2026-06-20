import http from "k6/http";

export const options = {
    vus: 1,                 //1 virtual used
    iterations: 150,            //runs 150 times, sends 150 requestsz
  };

export default function () {
	http.get("http://localhost:8080/api/data");
}