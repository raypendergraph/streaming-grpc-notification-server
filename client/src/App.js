import {useEffect, useState} from "react"
import './App.css';
import notifications from './generated/proto/notifications_pb';
import pb from './generated/proto/notifications_grpc_web_pb';
import { v4 as uuidv4 } from 'uuid';
import 'bulma/css/bulma.min.css';

const client = new pb.NotificationServicePromiseClient("", undefined, undefined)
const appUuid = uuidv4();
console.log("UUID", appUuid)
function App() {
    const [chat, setChat] = useState([])
    function sendChat(message){
        (async () => {
            const request = new notifications.SendMessageRequest()
                .setMessage(message)
                .setSender(appUuid)

            try {
                await client.sendMessage(request, undefined)
            } catch (e) {
                console.error(e)
            }
        })()
    }

    useEffect(() => {
        const req = new notifications.SubscribeRequest().setId(appUuid,undefined )
        const stream = client.subscribe(req,undefined);

        stream.on('data', function(resp) {
            setChat( current => [...current, resp.toObject()])
        });

        stream.on('status', function(status) {
            console.log(status.code);
            console.log(status.details);
            console.log(status.metadata);
        });

        stream.on('end', function(end) {
            setChat([])
            console.error("ended", end)
        });
    })
    console.log(chat.length)
    return (
        <section className="hero is-fullheight">
          <div className="hero-head">
            <header className="hero is-link is-bold">
              <div className="hero-body">
                <div className="container">
                  <p className="title">
                    Chat Messenger
                  </p>
                </div>
              </div>
            </header>
          </div>

          <div className="hero-body">
            <Messages chat={chat} />
          </div>

          <div className="hero-foot">
            <footer className="section is-small">
              <Chat saveMsg={sendChat} />
            </footer>
          </div>
        </section>
    )
}

const Chat = ({ saveMsg }) => (
    <form onSubmit={(e) => {
      e.preventDefault();
      saveMsg(e.target.elements.userInput.value);
      e.target.reset();
    }}>
      <div className="field has-addons">
        <div className="control is-expanded">
          <input className="input" name="userInput" type="text" placeholder="Type your message" />
        </div>
        <div className="control">
          <button className="button is-info">
            Send
          </button>
        </div>
      </div>
    </form>
);

const Messages = ({ chat }) => (

    <div style={{ height: '100%', width: '100%' }}>
      {chat.map((m, i) => {
        const msgClass = m.sender === appUuid ? 0 : 1
        return (
            <p key = {i} style={{ padding: '.25em', textAlign: msgClass ? 'left' : 'right', overflowWrap: 'normal' }}>
              <span className={`tag is-medium ${msgClass ? 'is-success' : 'is-info'}`}>{`${m.sender.substring(32)}: ${m.message}`}</span>
            </p>
        )}
      )}
    </div>
);


export default App