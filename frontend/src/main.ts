import {createApp} from 'vue'
import App from './App.vue'
import VueDatePicker from '@vuepic/vue-datepicker';
import './style.css';
import "bootstrap/dist/css/bootstrap.min.css"
import "bootstrap"
import '@vuepic/vue-datepicker/dist/main.css'

createApp(App).component('VueDatePicker', VueDatePicker).mount('#app');
