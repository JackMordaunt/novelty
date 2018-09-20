<template>
    <div>
        <pre v-for="(ii, status) in updates" :key="ii">
            {{JSON.stringify(status)}}
        </pre>
    </div>
</template>

<script>
// Loading component shows the user the status of an open media resource.
// This component cares about status updates.
export default {
    props: {
        show: Object,
    },
    created() {
        this.$ws.on("show.opened", (msg) => {
            // eslint-disable-next-line
            console.log(`player opened: ${msg}`)
        })
        this.$ws.on("show.status", (msg) => {
            this.updates.push(msg)
        })
        this.$ws.send("show.open", {
            name: this.show.name,
            uri: this.show.uri,
        })
    },
    data() {
        return {
            updates: [],
        }
    },
    methods: {
        cancel() {
            this.$ws.send("show.close", {
                name: this.show.name,
            })
        }
    },
}
</script>

<style>
</style>