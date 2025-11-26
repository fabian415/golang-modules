import { createRouter, createWebHistory } from 'vue-router';
import SampleBuildProgress from '../components/SampleBuildProgress.vue';

const routes = [
  {
    path: '/',
    name: 'SampleBuildProgress',
    component: SampleBuildProgress,
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;

