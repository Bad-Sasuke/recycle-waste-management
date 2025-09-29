import { defineStore } from 'pinia';
import { getCookie, deleteCookie } from '../stores/cookie'

/**
 * User interface - Expected to be returned by the /api/user/profile endpoint
 * The backend should include the 'role' field with values like 'admin', 'moderator', or 'user'
 */
interface User {
  user_id: string;
  username: string;
  email: string;
  image_url: string;
  created_at: string;
  last_login: string;
  role?: string; // Expected values: 'admin', 'moderator', 'user', etc.
}

export const useUsersStore = defineStore('users', {
    state: () => ({
        isLogin: false,
        jwt: '',
        user: null as User | null,
    }),
    getters: {
        profileImage: (state) => state.user?.image_url || '',
        username: (state) => state.user?.username || '',
    },
    actions: {
        async checkLogin() {
            const token = getCookie('token');
            this.isLogin = !!token;
            this.jwt = token || '';
            
            // If logged in, fetch user profile
            if (this.isLogin) {
                await this.fetchUserProfile();
            }
        },
        
        async fetchUserProfile() {
            if (!this.isLogin || !this.jwt) return;
            
            try {
                const apiUrl = import.meta.env.VITE_WEB_API;
                if (!apiUrl) {
                    console.error('API configuration error. Please check environment settings.');
                    return;
                }
                
                const response = await fetch(`${apiUrl}/api/user/profile`, {
                    headers: {
                        'Authorization': `Bearer ${this.jwt}`,
                        'Content-Type': 'application/json'
                    }
                });
                
                if (response.ok) {
                    const result = await response.json();
                    this.user = result.data;
                } else {
                    console.error(`Failed to fetch profile data: ${response.status} ${response.statusText}`);
                    // If profile fetch fails, set user to null to avoid stale data
                    this.user = null;
                }
            } catch (error) {
                console.error('Error fetching user profile:', error);
                this.user = null;
            }
        },
        
        logout() {
            deleteCookie('token');
            this.isLogin = false;
            this.jwt = '';
            this.user = null;
        }
    },
});