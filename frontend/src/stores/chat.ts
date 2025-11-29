import { defineStore } from 'pinia'
import { useUsersStore } from './users'
import { getCookie } from './cookie'

export interface ChatMessage {
  type: 'message' | 'join' | 'leave'
  customer_request_id: string
  sender_id: string
  sender_type: 'customer' | 'shop'
  message: string
  timestamp: string
}

export const useChatStore = defineStore('chat', {
  state: () => ({
    socket: null as WebSocket | null,
    isConnected: false,
    isPartnerOnline: false,
    messages: [] as ChatMessage[],
    currentRequestId: null as string | null,
    error: null as string | null,
  }),

  actions: {
    connect(requestId: string) {
      if (this.socket && this.isConnected && this.currentRequestId === requestId) {
        return // Already connected to the same room
      }

      // If switching rooms, close old connection but keep messages
      if (this.socket && this.currentRequestId !== requestId) {
        if (this.socket) {
          this.socket.close()
          this.socket = null
        }
        this.isConnected = false
        this.isPartnerOnline = false
        // Clear messages only when switching to different room
        this.messages = []
      }

      const userStore = useUsersStore()
      const user = userStore.user
      const token = getCookie('token')
      const apiUrl = import.meta.env.VITE_WEB_API || 'http://localhost:8080'
      const wsUrl = apiUrl.replace(/^http/, 'ws')

      if (!user || !token) {
        this.error = 'User not authenticated'
        return
      }

      this.currentRequestId = requestId
      // Reset online status when connecting
      this.isPartnerOnline = false

      // Construct WebSocket URL
      const url = `${wsUrl}/ws/chat?user_id=${user.user_id}&request_id=${requestId}`

      console.log('Connecting to WebSocket:', url)

      try {
        this.socket = new WebSocket(url)

        this.socket.onopen = () => {
          console.log('WebSocket connected (onopen)')
          this.isConnected = true
          this.error = null
        }

        this.socket.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data) as ChatMessage

            // Handle presence events
            if (message.type === 'join') {
              if (message.sender_id !== userStore.user?.user_id) {
                this.isPartnerOnline = true
                console.log('Partner joined/is online')
              }
            } else if (message.type === 'leave') {
              if (message.sender_id !== userStore.user?.user_id) {
                this.isPartnerOnline = false
                console.log('Partner left')
              }
            }

            // Improved duplicate detection
            // Check if there's a message with same sender_id and message content sent within 5 seconds
            const now = new Date(message.timestamp).getTime()
            const isDuplicate = this.messages.some((m) => {
              const msgTime = new Date(m.timestamp).getTime()
              const timeDiff = Math.abs(now - msgTime)

              return (
                m.sender_id === message.sender_id &&
                m.message === message.message &&
                timeDiff < 5000 // Within 5 seconds
              )
            })

            if (!isDuplicate) {
              this.messages.push(message)
              console.log('Received message:', message)
            } else {
              console.log('Duplicate message ignored:', message)
            }
          } catch (e) {
            console.error('Failed to parse message:', e)
          }
        }

        this.socket.onclose = (event) => {
          console.log('WebSocket disconnected (onclose)', event.code, event.reason)
          this.isConnected = false
          this.isPartnerOnline = false
          this.socket = null
        }

        this.socket.onerror = (error) => {
          console.error('WebSocket error:', error)
          this.error = 'Connection error'
          this.isConnected = false
        }
      } catch (e: unknown) {
        console.error('Failed to create WebSocket connection:', e)
        if (e instanceof Error) {
          this.error = e.message
        } else {
          this.error = 'Unknown error occurred'
        }
      }
    },

    disconnect() {
      console.log('Disconnecting WebSocket...')
      if (this.socket) {
        this.socket.close()
        this.socket = null
      }
      this.isConnected = false
      // Don't clear messages - keep chat history
      // this.messages = []
      this.currentRequestId = null
    },

    async sendMessage(message: string) {
      console.log('sendMessage called with:', message)

      if (!this.socket) {
        console.error('Cannot send message: Socket is null')
        // Try to reconnect if we have request ID
        if (this.currentRequestId) {
          console.log('Attempting to reconnect...')
          this.connect(this.currentRequestId)
          // Wait a bit for connection
          await new Promise((resolve) => setTimeout(resolve, 1000))
        } else {
          return
        }
      }

      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        const userStore = useUsersStore()

        // Optimistic update - add message immediately
        const optimisticMessage: ChatMessage = {
          type: 'message',
          customer_request_id: this.currentRequestId || '',
          sender_id: userStore.user?.user_id || '',
          sender_type: userStore.user?.role === 'moderator' ? 'shop' : 'customer',
          message: message,
          timestamp: new Date().toISOString(),
        }

        // Add to local messages immediately for instant UI update
        this.messages.push(optimisticMessage)

        const payload = {
          customer_request_id: this.currentRequestId,
          message: message,
        }
        console.log('Sending payload:', payload)
        try {
          this.socket.send(JSON.stringify(payload))
          console.log('Message sent successfully')
        } catch (e) {
          console.error('Error sending message:', e)
          // Remove optimistic message if send failed
          const index = this.messages.indexOf(optimisticMessage)
          if (index > -1) {
            this.messages.splice(index, 1)
          }
        }
      } else {
        console.error(
          'Cannot send message: Socket is not OPEN. State:',
          this.socket ? this.socket.readyState : 'null',
        )
      }
    },
  },
})
