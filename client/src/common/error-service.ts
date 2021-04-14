import { message } from 'antd';

export const handleError = (e: any) => {
  console.error(e);
  message.error('Error happened.');
}

export const handleFetchError = (r: Response): Response => {
  if (!r.ok) {
    throw Error(r.statusText);
  }
  return r;
}
