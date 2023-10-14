<template>
  <form @click.prevent="onSubmit">
    <div v-if="sessionUser">
      <div class="chat-history" ref="chatHistory">
        <ul>
          <li v-for="post in posts">
            <span class="chat-history__user">{{ post.user.username }}</span> :
            <span class="chat-history__message">{{ post.message }}</span>
            <span class="chat-history__timestamp">{{ post.timestamp }}</span>
          </li>
        </ul>
      </div>
      <div class="chat-input">
      <input class="message" v-model="message" type="text" >
      <input class="button" type="submit" value="Send" @click="sendMessage">
        <input class="button" type="submit" value="Logout" @click="logout">
      </div>
    </div>
    <div v-else>
      <input class="user" type="text" v-model="username" placeholder="Username">
      <input class="password" type="password" v-model="password" placeholder="Password">
      <input class="button" type="submit" value="Login" @click="login">
      <div class="alert" v-if="!userValid">
        Invalid username or password. Try again!
      </div>
    </div>
  </form>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      message: "",
      socket: null,
      posts: [],
      username: "UserOne",
      password: "12345",
      sessionUser: "",
      userValid : true,
    }
  },

  methods: {
    instanceSocket() {
      this.socket = new WebSocket("ws://localhost:5000/ws")

      this.socket.onmessage = (msg) => {
        this.acceptMsg(msg)
      }

      this.socket.onopen = (evt) => {
        let msg = {
          userID: this.sessionUser.id,
          user: {
            id: this.sessionUser.id,
            username: this.sessionUser.username,
          },
          message: "<SayHi>"
        }
        this.socket.send(JSON.stringify(msg))
      }
    },

    sendMessage() {
      if(this.message.trim() === "") {
        return
      }

      let msg = {
        userID: this.sessionUser.id,
        user: {
          id: this.sessionUser.id,
          username: this.sessionUser.username,
        },
        message: this.message
      }
      this.socket.send(JSON.stringify(msg))
      this.message = ''
    },

    acceptMsg(msg) {
      this.posts = JSON.parse(msg.data).reverse().map(p => {
        if(p === null || p.timestamp === null || p.timestamp === undefined) {
          return p
        }

        const date = new Date(p.timestamp)
        p.timestamp = date.toLocaleDateString('en-US', { weekday: 'long', hour: "numeric", minute: "numeric" })

        return p
      })
    },

    async login() {
      let user = {
        username: this.username,
        password: this.password,
      }

      const res = await fetch("http://localhost:5000/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(user),
      })

      res.json().then((user) => {
        if(!user.username) {
          this.userValid = false
        } else {
          sessionStorage.user = JSON.stringify(user)
          this.sessionUser = user
          this.userValid = true

          this.instanceSocket()
        }
      })
    },

    logout() {
      let msg = {
        userID: this.sessionUser.id,
        user: {
          id: this.sessionUser.id,
          username: this.sessionUser.username,
        },
        message: "<SayBye>"
      }
      this.socket.send(JSON.stringify(msg))

      delete(sessionStorage.user)
      this.sessionUser = null
      this.userValid = true
      this.username = ""
      this.password = ""
    }
  },
}
</script>

<style>
#app {
 font-family: sans-serif;
}

.chat-input {
  margin-top: 20px;
}

.chat-history {
  width: 93%;
  height: 500px;
  border: black 1px solid;
  background: white;
  color: black;
  overflow-y: scroll;
}

.chat-history__user {
  color: blueviolet;
}

.chat-history > ul > li {
  list-style-type: none;
  margin: 10px;

  border-bottom: 1px solid lightgray;

}

.chat-history__timestamp {
  color: gray;
  float: right;
  font-size: small;
}


.message {
  width: 70%;
}

.alert {
  color: red;
}

.button {
  width: 10%;
  margin-left: 20px;
}
</style>