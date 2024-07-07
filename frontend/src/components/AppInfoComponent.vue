<script lang="ts" setup>
import {reactive} from 'vue'
import {BackupFolder} from '../../wailsjs/go/app/App'
import {player} from "../../wailsjs/go/models";

interface Props {
  selectedPlayer: player.Player
}

const props = defineProps<Props>()

interface Data {
  backupFolder?: string
}

const data = reactive<Data>({
  backupFolder: undefined,
})

BackupFolder(props.selectedPlayer).then(result => data.backupFolder = result)
</script>

<template>
  <div class="row" v-if="data.backupFolder">
    <div class="col">
      <div class="alert alert-dark fade show" role="alert">
        <h4 class="alert-heading fw-bold">Hi {{ props.selectedPlayer.battleTag }}!</h4>
        <p class="mb-0">This application will let you restore your game history to any checkpoint you would want to.</p>
        <p class="mb-0">Upon restoring the application will automatically create backups of your history before
          executing any operations.</p>
        <p>in case of an error you can restore.</p>
        <hr>
        <p class="mb-0"><strong>The backup location for your account is:</strong><br/><u>{{ data.backupFolder }}</u></p>
      </div>
    </div>
  </div>
</template>
