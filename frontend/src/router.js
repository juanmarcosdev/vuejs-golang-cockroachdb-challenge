import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/views/Home'
import Endpoint1 from '@/views/Endpoint1'
import Endpoint2 from '@/views/Endpoint2'
import Error from '@/views/Error'

Vue.use(Router)

export default new Router({
    mode: 'history',
    routes: [
        {
            path: '/',
            name: 'home',
            component: Home
        },
        {
            path: '/endpoint1/post',
            name: 'endpoint1',
            component: Endpoint1
        },
        {
            path: '/endpoint2/get',
            name: 'endpoint2',
            component: Endpoint2
        },
        {
            path: '*',
            name: 'error',
            component: Error
        }
    ]
})