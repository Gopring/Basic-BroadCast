import { createPeerConnection, makeRequestBody, fetchFromServer, getUserId, videoElement } from './utils.js';

const viewServerUrl="http://localhost:8080/channel/view"

export const startView = async () => {
    try {
        const key = getUserId();

        // Create a new peer connection to connect to the server
        const pc = createPeerConnection();
        pc.addTransceiver('video')

        // Create an offer to receive the stream from the server
        const offer = await pc.createOffer();
        await pc.setLocalDescription(offer);

        pc.ontrack = function (event) {
            const el = videoElement
            el.srcObject = event.streams[0]
            el.autoplay = true
            el.controls = true
        }

        // Send the offer to the server to start viewing
        const requestBody = makeRequestBody(key, offer.sdp);
        const res = await fetchFromServer(viewServerUrl, requestBody, key);
        const sdpAnswer =await res.text()
        // Set the server's SDP answer to establish connection
        const remoteDescription = new RTCSessionDescription({
            type: 'answer',
            sdp: sdpAnswer
        });
        await pc.setRemoteDescription(remoteDescription);
    } catch (error) {
        console.error('Error viewing:', error);
    }
};
