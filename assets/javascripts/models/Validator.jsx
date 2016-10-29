/* eslint max-len: ["off"] */
/**
 * Created by shudi on 2016/10/22.
 */
export default {
  isMobile(mobile) {
    return /^1[3|4|5|7|8][0-9]{9}$/.test(mobile);
  },

  isEmail(email) {
    return /^([a-z0-9A-Z]+[-|\\.]?)+[a-z0-9A-Z]@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)?\\.)+[a-zA-Z]{2,}$/.test(email);
  },

  isStudentId(studentId) {
    return /^\d{10}$/.test(studentId);
  },
};
