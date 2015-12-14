import time
import sys
import shutil
import os

from flask import *

app = Flask(__name__)
app.debug = True


playlist = ['short1.mp3','short2.mp3','long2.mp3','long3.mp3','short3.mp3']
current_song = -1
last_activated = 0
is_playing = False

def getTime():
    return int(time.time()*1000)


@app.route("/")
def index_html():
    return render_template('index.html')

@app.route("/sync", methods=['GET', 'POST'])
def sync():
    #searchword = request.args.get('key', '')
    if request.method == 'POST':
        print(getTime())
        data = {}
        data['client_timestamp'] = request.form['client_timestamp']
        data['server_timestamp'] = getTime()
        data['next_song'] = next_song_time
        data['is_playing'] = is_playing
        return jsonify(data)


@app.route("/nextsong", methods=['GET', 'POST'])
def finished():
    response = {'message':'loading!'}
    if request.method == 'POST':
        nextSong(6)
    return jsonify(response)

@app.route("/playing", methods=['GET', 'POST'])
def playing():
    global is_playing
    response = {'message':'loading!'}
    if request.method == 'POST':
        is_playing = True
    return jsonify(response)

def nextSong(delay):
    global last_activated
    global current_song
    global next_song_time
    global is_playing
    if time.time() - last_activated > 10:
        is_playing = False
        current_song += 1
        last_activated = time.time()
        shutil.copy('./' + playlist[current_song],'./static/')
        os.rename('./static/' + playlist[current_song],'./static/sound.mp3')
        os.system('scp ' + playlist[current_song] + ' phi@192.168.1.11:/www/data/sound.mp3')
        next_song_time = getTime() + delay*1000
        print ('next up: ' + playlist[current_song])

if __name__ == "__main__":
    nextSong(20)
    app.run(host='192.168.1.2')

