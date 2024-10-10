const videoElement = document.getElementById('videoElement');
const startBroadcastButton = document.getElementById('startBroadcast');
const startViewButton = document.getElementById('startView');
const userIdInput = document.getElementById('userIdInput');

let localStream;
let peerConnection;

const broadcastServerUrl = 'http://localhost:8080/channel/broadcast';
const viewServerUrl = 'http://localhost:8080/channel/view';

const getUserId = () => {
  const key = userIdInput.value;
  if (!key) {
    alert('Please enter an ID');
    throw new Error('User ID is required');
  }
  return key;
};

const createPeerConnection = (onTrackCallback) => {
  const config = {
    iceServers: [
      {
        urls: 'stun:stun.l.google.com:19302'
      }
    ]
  };
  const pc = new RTCPeerConnection(config);

  pc.onicecandidate = (event) => {
    if (event.candidate) {
      console.log('New ICE Candidate:', event.candidate);
    }
  };

  if (onTrackCallback) {
    pc.ontrack = onTrackCallback;
  }

  return pc;
};

const addLocalTracksToPeerConnection = (pc, stream) => {
  stream.getTracks().forEach(track => {
    pc.addTrack(track, stream);
  });
};

const makeRequestBody = (id, sdp) => {
  return JSON.stringify({
    id: id,
    sdp: sdp
  });
};

const fetchFromServer = (url, requestBody, apiKey) => {
  return fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'api-key': apiKey
    },
    body: requestBody
  }).then(response => response.json());
};

const startBroadcast = async () => {
  const key = getUserId();

  localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true });
  videoElement.srcObject = localStream;

  peerConnection = createPeerConnection();
  addLocalTracksToPeerConnection(peerConnection, localStream);

  const offer = await peerConnection.createOffer();
  await peerConnection.setLocalDescription(offer);

  const requestBody = makeRequestBody(key, offer.sdp);
  fetchFromServer(broadcastServerUrl, requestBody, key)
      .then(async data => {
        const remoteDescription = new RTCSessionDescription({
          type: 'answer',
          sdp: data.sdp
        });
        await peerConnection.setRemoteDescription(remoteDescription);
      }).catch(error => {
    console.error('Error broadcasting:', error);
  });
};

const startView = async () => {
  const key = getUserId();

  peerConnection = createPeerConnection((event) => {
    videoElement.srcObject = event.streams[0];
  });

  const requestBody = makeRequestBody(key, '');
  fetchFromServer(viewServerUrl, requestBody, key)
      .then(async data => {
        const remoteDescription = new RTCSessionDescription({
          type: 'offer',
          sdp: data.sdp
        });
        await peerConnection.setRemoteDescription(remoteDescription);

        const answer = await peerConnection.createAnswer();
        await peerConnection.setLocalDescription(answer);

        const answerRequestBody = makeRequestBody(key, answer.sdp);
        fetchFromServer(viewServerUrl, answerRequestBody, key);
      }).catch(error => {
    console.error('Error viewing:', error);
  });
};

startBroadcastButton.onclick = startBroadcast;
startViewButton.onclick = startView;