<script lang="ts" setup>
import {history} from "../../wailsjs/go/models";
import {reactive, ref} from "vue";

interface Props {
  game: history.GameHistory
}

const props = defineProps<Props>()

interface Data {
  toggled: boolean,
}

const data = reactive<Data>({
  toggled: false
})

function toggleStats() {
  data.toggled = !data.toggled
}
</script>

<template>
  <tr v-on:click="toggleStats">
    <td class="text-center">{{ props.game.gameId }}</td>
    <td>{{ props.game.date }}</td>
    <td>
      {{ props.game.elo }}
      <span v-if="['WIN', 'LOSS'].includes(props.game.outcome)" class="badge badge rounded-pill"
            :class="props.game.outcome == 'WIN' ? 'bg-success' : 'bg-danger'">
      {{ `${props.game.outcome == 'WIN' ? '+' : '-'}${props.game.eloDiff}` }}
    </span>
    </td>
    <td>
      {{ props.game.outcome }}
      <span v-if="props.game.winsStreak > 1">&#128293; {{ props.game.winsStreak }}</span>
    </td>
  </tr>
</template>
