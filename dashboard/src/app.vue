<style scoped>
    @import 'styles/common.css';
</style>
<template>
    <div>
        <router-view></router-view>
    </div>
</template>
<script>
    export default {
        data () {
            return {}
        },
        mounted () {

        },
        beforeMount() {
            this.axios.defaults.baseURL = 'http://localhost:20101';
        },
        beforeDestroy () {

        },
        methods: {

        }
    }
</script>
