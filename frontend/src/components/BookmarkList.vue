<template>
  <div>
    <div class="block control">
      <input class="input" type="search" placeholder="Search" v-model="filter" autofocus @search="onSearch">
    </div>
    <div class="block bookmark" v-for="bookmark in bookmarks" :key="bookmark.ID">
      <p class="has-text-weight-bold">{{ bookmark.Title }}</p>
      <p class="is-size-7"><a :href="bookmark.URL">{{ bookmark.URL }}</a></p>
      <p class="content">{{ bookmark.Content|excerpt }}</p>
      <p class="block has-text-right">
        <a @click.prevent="onRemoveClicked(bookmark)" class="button is-small is-danger is-outlined">Remove</a>
        <a @click.prevent="onReadItLaterClicked(bookmark)" class="button is-small is-primary" v-if="bookmark.Archived">Read it later</a>
        <a @click.prevent="onArchiveClicked(bookmark)" class="button is-small is-dark" v-else="bookmark.Archived">Archive</a>
      </p>
    </div>
  </div>
</template>

<script>
import { mapActions } from 'vuex'

export default {
  data () {
    return {}
  },
  computed: {
    filter: {
      get () {
        return this.$store.state.bookmarks.filter
      },
      set (value) {
        this.$store.commit('filter', value)
      }
    },
    bookmarks () {
      return this.$store.state.bookmarks.bookmarks
    }
  },
  filters: {
    excerpt: function (value) {
      if (!value) {
        return 'No content'
      }
      return value.toString().substring(0, 500) + '...'
    }
  },
  methods: {
    onSearch (event) {
      this.$router.push({query: {q: this.filter}})
    },
    ...mapActions({
      onReadItLaterClicked: 'readLaterBookmark',
      onArchiveClicked: 'archiveBookmark',
      onRemoveClicked: 'removeBookmark'
    })
  }
}
</script>

<style>
.bookmark {
  border-bottom: 1px solid hsl(0, 0%, 96%);
  padding-bottom: 1.5rem;
}
.content:not(:last-child) {
  margin-bottom: 0.5rem;
}
</style>
