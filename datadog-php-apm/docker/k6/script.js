import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 1, // ユーザー数を1に設定
  duration: '100m', // durationを無期限（0）に設定
};

export default function () {
  while (true) {
    const res = http.get('http://nginx:80/db.php');
    check(res, {
      'status was 200': (r) => r.status === 200,
    });
    sleep(10); // 1秒間隔でリクエストを送信
  }
}