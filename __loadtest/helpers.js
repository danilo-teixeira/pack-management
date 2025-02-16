export function generateRandomString(length) {
  const charset = 'abcdefghijklmnopqrstuvwxyz';
  let res = '';
  while (length--) res += charset[Math.random() * charset.length | 0];
  return res;
}

export function generateRandomDate() {
  const start = new Date(2021, 0, 1);
  const end = new Date(2025, 0, 1);
  const date = new Date(start.getTime() + Math.random() * (end.getTime() - start.getTime()));

  let month = date.getMonth() + 1;
  if (month < 10) {
    month = '0' + month
  }

  let day = date.getDate();
  if (day < 10) {
    day = '0' + day
  }

  return date.getFullYear() + '-' + month + '-' + day;
}
