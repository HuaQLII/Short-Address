import axios from "axios";
axios.defaults.baseURL = "http://localhost:9090/api/v1";
export function longAdd(data) {
    return axios.request({
        url: '/long?url="asdasdsad"',
        method: "post",
        data:data
    });
}



// axios.post('long')
//     .then(res => {
//         console.log(res)
//     })
//     .catch(err => {
//         console.error(err);
//     })