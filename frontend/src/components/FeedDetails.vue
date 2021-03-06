<template>
  <div class="block">
    <div class="pull-right">
      <a class="button is-small is-primary is-outlined" @click.prevent="onRefreshFeedClicked(feed)">Refresh Feed</a>
      <a class="button is-small is-danger is-outlined" @click.prevent="onDeleteFeedClicked(feed)">Delete Feed</a>
    </div>
    <h1 class="title">{{ feed.Title }}</h1>
    <h2 class="subtitle">{{ feed.URL }}</h2>
    <p><i :title="feed.LastAuthored|moment('dddd, MMMM Do YYYY, HH:mm')">Last item created {{ feed.LastAuthored|moment("from") }}</i></p>
    <p><i :title="feed.Refreshed|moment('dddd, MMMM Do YYYY, HH:mm')">Last refreshed {{ feed.Refreshed|moment("from") }}</i></p>
    <hr />
    <div v-if="items.length > 0">
      <div class="block feed-item" v-for="item in items" :key="item.ID">
        <p class="has-text-weight-bold">{{ item.Title }}</p>
        <p class="is-size-7"><a :href="item.URL">{{ item.URL }}</a></p>
        <p class="content">{{ item.Content|excerpt }}</p>
        <p class="block has-text-right">
          <a @click.prevent="onRemoveClicked(item)" class="button is-small is-danger is-outlined">Remove</a>
          <a @click.prevent="onReadItLaterClicked(item)" class="button is-small is-primary">Read it later</a>
        </p>
      </div>
    </div>
    <div v-else class="block has-text-centered">
      <i>No items in this feed</i>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  computed: {
    feed () {
      return this.$store.getters.selectedFeed
    },
    items () {
      return this.$store.getters.items
    }
  },
  methods: {
    ...mapActions({
      onRefreshFeedClicked: 'refreshFeed',
      onDeleteFeedClicked: 'deleteFeed',
      onReadItLaterClicked: 'readLaterFeedItem',
      onRemoveClicked: 'removeFeedItem'
    })
  },
  filters: {
    excerpt (value) {
      if (!value) {
        return 'No content'
      }
      return value.toString().substring(0, 500) + '...'
    }
  }
}
</script>

<style>
.feed-item {
  border-bottom: 1px solid hsl(0, 0%, 96%);
  padding-bottom: 1.5rem;
}
.content:not(:last-child) {
  margin-bottom: 0.5rem;
}
</style>
