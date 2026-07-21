import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../types/auth';

interface AuthState {
    user: User | null;
    accessToken: string | null;
    refreshToken: string | null;
    isAuthenticated: boolean;
    loading: boolean;
    setUser: (user: User | null) => void;
    setTokens: (tokens: { accessToken: string; refreshToken: string }) => void;
    login: (user: User, tokens: { accessToken: string; refreshToken: string }) => void;
    logout: () => void;
    setLoading: (loading: boolean) => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            accessToken: null,
            refreshToken: null,
            isAuthenticated: false,
            loading: false,

            setUser: (user) => set({ user, isAuthenticated: !!user }),
            
            setTokens: (tokens) => set({ 
                accessToken: tokens.accessToken, 
                refreshToken: tokens.refreshToken 
            }),
            
            login: (user, tokens) => set({ 
                user, 
                accessToken: tokens.accessToken, 
                refreshToken: tokens.refreshToken, 
                isAuthenticated: true 
            }),
            
            logout: () => set({ 
                user: null, 
                accessToken: null, 
                refreshToken: null, 
                isAuthenticated: false 
            }),
            
            setLoading: (loading) => set({ loading }),
        }), // Fixed: Cleaned up structural closures here
        {
            name: 'auth-storage',
            partialize: (state) => ({
                accessToken: state.accessToken,
                refreshToken: state.refreshToken,
                user: state.user,
                isAuthenticated: state.isAuthenticated,
            }),
        }
    )
);
