import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import MainLayout from '../layouts/MainLayout.vue';
import Login from '../views/Login.vue';
import Users from '../views/Users.vue';
import UserGroups from '../views/UserGroups.vue';
import FileBrowser from '../views/FileBrowser.vue';
import FileUpload from '../views/FileUpload.vue';

const routes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: Login },
  {
    path: '/',
    component: MainLayout,
    children: [
      { path: '', redirect: '/users' },
      { path: 'users', name: 'users', component: Users, meta: { requiresAuth: true } },
      { path: 'user-groups', name: 'user-groups', component: UserGroups, meta: { requiresAuth: true } },
      { path: 'file-browser', name: 'file-browser', component: FileBrowser, meta: { requiresAuth: true } },
      { path: 'file-upload', name: 'file-upload', component: FileUpload, meta: { requiresAuth: true } },
      { path: 'audit-logs', name: 'audit-logs', component: FileBrowser, meta: { requiresAuth: true } },
      { path: 'settings', name: 'settings', component: FileBrowser, meta: { requiresAuth: true } },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token') || '';
  if (to.meta?.requiresAuth && !token) {
    next({ name: 'login', query: { redirect: to.fullPath } });
  } else {
    next();
  }
});

export default router;