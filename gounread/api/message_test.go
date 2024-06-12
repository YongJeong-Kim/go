package api

/*func TestSendMessage(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/rooms/asdasd/send", nil)
	require.NoError(t, err)
	defer ws.Close()

	err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	require.NoError(t, err)
	_, p, err := ws.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, string(p), "hello")
}*/
