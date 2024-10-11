export const videoElement = document.getElementById('videoElement');
export const userIdInput = document.getElementById('userIdInput');

export const getUserId = () => {
    const key = userIdInput.value;
    if (!key) {
        alert('Please enter an ID');
        throw new Error('User ID is required');
    }
    return key;
};

export const createPeerConnection = () => {
    const config = {
        iceServers: [
            {
                urls: 'stun:stun.l.google.com:19302'
            }
        ]
    };
    return new RTCPeerConnection(config);
};

export const addLocalTracksToPeerConnection = (pc, stream) => {
    stream.getTracks().forEach(track => {
        pc.addTrack(track, stream);
    });
};

export const makeRequestBody = (id, sdp) => {
    return JSON.stringify({
        id: id,
        sdp: sdp
    });
};

export const fetchFromServer = async (url, requestBody, apiKey) => {
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'api-key': apiKey
            },
            body: requestBody
        });
        return await response;
    } catch (error) {
        console.error('Error fetching from server:', error);
        throw error;
    }
};
