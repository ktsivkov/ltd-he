<script lang="ts" setup>
import {EventsOn} from '../../wailsjs/runtime'
import {reactive} from "vue";

interface Event {
  type: string,
  message: string,
}

interface Data {
  events: Array<Event>,
}

const data = reactive<Data>({
  events: [],
})

const getAlertBg = (e: Event): string => {
  switch (e.type){
    case "success":
      return "text-bg-success"
    case "error":
      return "text-bg-danger"
    case "warning":
      return "text-bg-warning"
    case "info":
      return "text-bg-info"
    default:
      console.error(`unknown event type ${e.type}`)
      return "text-bg-primary"
  }
}


const onAlertEvent = (e: Event) => {
  data.events.push(e)
}

EventsOn("alert", onAlertEvent)

</script>

<template>
  <div aria-live="polite" aria-atomic="true" class="position-relative" v-if="data.events">
    <div class="toast-container top-0 end-0 p-3">
      <div v-for="event in data.events" class="toast align-items-center border-0 fade show" v-bind:class="getAlertBg(event)" role="alert" aria-live="assertive" aria-atomic="true">
        <div class="d-flex">
          <div class="toast-body"><pre style="white-space: pre-wrap;">{{event.message}}</pre></div>
          <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
        </div>
      </div>
    </div>
  </div>
</template>
