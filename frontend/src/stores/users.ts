import { defineStore } from 'pinia';
import { getCookie, deleteCookie } from '../stores/cookie'

export const useUsersStore = defineStore('users', {
    state: () => ({
        isLogin: false,
        jwt: '',
    }),
    getters: {},
    actions: {
        checkLogin() {
            const token = getCookie('token');
            this.isLogin = !!token;
            this.jwt = token || '';
        },
        logout() {
            deleteCookie('token');
            this.isLogin = false;
            this.jwt = '';
        }
    },
});