import { createApp } from 'vue';
import 'ant-design-vue/dist/reset.css';

import App from './App.vue';
import { router } from './router';
import { pinia } from './stores/auth';
import './styles/base.css';

const app = createApp(App);

app.use(pinia);
app.use(router);

app.mount('#app');
