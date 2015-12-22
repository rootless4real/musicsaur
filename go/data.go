package main

import "sort"

var conf tomlConfig
var statevar State
var rawSongData []byte

// Data for configuration file

type tomlConfig struct {
	ClientData       clientInfo  `toml:"raspberry_pis"`
	ClientParameters clientParms `toml:"client_parameters"`
	ServerParameters serverParms `toml:"server_parameters"`
}

type clientInfo struct {
	Clients string `toml:"clients"`
}

type clientParms struct {
	CheckUpWaitTime int `toml:"check_up_wait_time"`
	MaxSyncLag      int `toml:"max_sync_lag"`
}

type serverParms struct {
	MusicFolder         string `toml:"music_folder"`
	Port                int    `toml:"port"`
	TimeToNextSong      int    `toml:"time_to_next_song"`
	TimeToDisallowSkips int    `toml:"time_to_disallow_skips"`
}

// Data for state

type State struct {
	SongMap          map[string]Song
	SongList         sort.StringSlice
	PathList         map[string]bool
	SongStartTime    int64
	IsPlaying        bool
	CurrentSong      string
	CurrentSongIndex int
}

// Data for Song

type SyncJSON struct {
	Current_song     string  `json:"current_song"`
	Client_timestamp int64   `json:"client_timestamp"`
	Server_timestamp int64   `json:"server_timestamp"`
	Is_playing       bool    `json:"is_playing"`
	Song_time        float64 `json:"song_time"`
	Song_start_time  int64   `json:"next_song"`
}

type Song struct {
	Fullname string
	Title    string
	Artist   string
	Album    string
	Path     string
	Length   int64
}

type IndexData struct {
	PlaylistHTML    string
	RandomInteger   int64
	CheckupWaitTime int64
	MaxSyncLag      int64
}

var index_html2 = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">

    <title>MusicSAUR</title>
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta content="no-cache" http-equiv="Cache-control">
    <meta content="-1" http-equiv="Expires">
    <script src="/math.js" type="text/javascript">
    </script>
    <script src="/jquery.js" type="text/javascript">
    </script>
    <script src="/howler.js" type="text/javascript">
    </script>
</head>
<body>
<audio controls preload="auto" src="./sound.mp3" id="sound" type="audio/mpeg">
  Your browser does not support the audio tag.
</audio>
</body>
</html>
`

var index_html = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">

    <title>MusicSAUR</title>
    <meta content="width=device-width, initial-scale=1" name="viewport">
    <meta content="no-cache" http-equiv="Cache-control">
    <meta content="-1" http-equiv="Expires">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/mathjs/2.5.0/math.min.js" type="text/javascript">
    </script>
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js">
    </script>
    <script src="/static/howler.js" type="text/javascript">
    </script>
    <!--
<audio preload="auto" src="/static/sound.mp3?{{ data['random_integer'] }}" id="sound" type="audio/mpeg">
  Your browser does not support the audio tag.
</audio>-->
    <link href="/static/normalize.css" rel="stylesheet">
    <link href="/static/skeleton.css" rel="stylesheet">
    <style>
    a { cursor: pointer; }

    .u-pull-right {
  float: left; }
    </style>
    <script>

var sound = new Howl({
  src: ['/sound.mp3?{{ data['random_integer'] }}'],
  preload: true
});

    </script>
</head>

<body>
<script>

var time = Date.now || function() {
return +new Date.getTime();
}

// CONSTANTS
var CHECK_UP_WAIT_TIME = {{ data['check_up_wait_time'] }};
var CHECK_UP_ITERATION = 1;
var check_up_counter = 0;
var MAX_SYNC_LAG = {{ data['max_sync_lag'] }};

// GLOBALS
var lagTimes = [];
var tryWait = 0;
var computeTimes = [];
var correct_time_delta = [];
var correct_latency = [];
var next_trigger = time() + 1000000;
var true_time_delta = 0;
var true_server_time_delta = 0;
var sound_activated = false;
var seconds_left = 0;
var current_song = "None"
var current_song_name = "None"
var secondTimeout3 = setTimeout(function() {
console.log('3 seconds left')
}, 100000);
var secondTimeout2 = setTimeout(function() {
console.log('2 seconds left')
}, 100000);
var secondTimeout1 = setTimeout(function() {
console.log('1 seconds left')
}, 100000);
var secondTimeout0 = setTimeout(function() {
console.log('0 seconds left')
}, 100000);
var mainInterval = 0;
var runningDiff = 0;



function makeRequests(callback) {

for (var i = 0; i < 23; i++) {

  setTimeout(function postRequest() {

// Send the data using post
    var posting = $.post('/sync', {
        'client_timestamp': time(),
        'current_song': current_song
    });

    // Put the results in a div
    posting.done(function(data) {
        var timeNow = time();
        current_song = data['current_song']
        latency = timeNow - data['client_timestamp']
        half_latency = latency / 2.0
        time_delta = timeNow - data['server_timestamp']
        next_trigger = data['next_song']

        correct_time_delta.push(time_delta + half_latency);
        correct_latency.push(half_latency);
        if (correct_time_delta.length==23) {
          console.log('correct_time_delta');
          console.log(correct_time_delta);
          var mean = math.mean(correct_time_delta);
          var median = math.median(correct_time_delta);
          var std = math.std(correct_time_delta);
          var sum = 0
          var num = 0
          for (var j = 0; j < correct_time_delta.length; j++) {
              if (correct_time_delta[j]<median+std) {
                  sum = sum + correct_time_delta[j];
                  num = num + 1;
              }
          }
          true_time_delta = sum / num;

          var mean = math.mean(correct_latency);
          var median = math.median(correct_latency);
          var std = math.std(correct_latency);
          var sum = 0
          var num = 0
          for (var j = 0; j < correct_latency.length; j++) {
              if (correct_latency[j]<median+std) {
                  sum = sum + correct_latency[j];
                  num = num + 1;
              }
          }
          true_server_time_delta = sum / num;

          clearTimeout(secondTimeout3);
          secondTimeout3 = setTimeout(function() {
              console.log('3 seconds left');
              $("div.info1").text('Playing in 3...');
          }, next_trigger - (time() - true_time_delta) - 3000);
          clearTimeout(secondTimeout2);
          secondTimeout2 = setTimeout(function() {
              console.log('2 seconds left');
              $("div.info1").text('Playing in 2...');
          }, next_trigger - (time() - true_time_delta) - 2000);
          clearTimeout(secondTimeout1);
          secondTimeout1 = setTimeout(function() {
              console.log('1 seconds left');
              $("div.info1").text('Playing in 1...');
          }, next_trigger - (time() - true_time_delta) - 1000);
          clearTimeout(secondTimeout0);
          secondTimeout0 = setTimeout(function() {
              console.log('playing song');
              current_song_name = current_song.split(":");
              current_song_name = current_song_name[current_song_name.length-1];
              $("div.info1").html('Loading <b>' + current_song_name + '</b>...');
              mainInterval = setInterval(function(){
                checkIfSkipped();
              }, CHECK_UP_WAIT_TIME);
              sound.play();
              if (data['is_playing']==true) {
                  sound.seek(data['song_time'])
              }
              // var posting = $.post('/playing', {
              // 'message': 'im playing a song'
              // });

              // // Put the results in a div
              // posting.done(function(data) {

              // });
          }, next_trigger - (time() - true_time_delta));

        }
    });

  }, i*180 );
    
}
}


function checkIfSkipped() {

    
    // Send the data using post
    var posting = $.post('/sync', {
        'client_timestamp': time(),
        'current_song': current_song
    });

    // Put the results in a div
    posting.done(function(data) {
        var start = new Date().getTime();
        var time_delta2 = time()-(data['server_timestamp']+true_time_delta);
      check_up_counter = check_up_counter + 1;
      if (data['is_playing']==false) {
        console.log('reloading page');
        sound.unload()
        location.reload(true);
      } else if (check_up_counter %% CHECK_UP_ITERATION==0) {
        check_up_counter = 0;
        var mySongTime = sound.seek();
        if (typeof(mySongTime)=="object") {
          mySongTime = 0;
          console.log('Still loading...')
            $("div.info1").html('Loading <b>' + current_song_name + '</b>...');
        }

        if (mySongTime == 0) {
          sound.seek(data['song_time']+time_delta2/1000.0);
        } else {
          var diff = data['song_time']+time_delta2/1000.0 - mySongTime;
          if (Math.abs(diff) > MAX_SYNC_LAG/1000.0) {
            CHECK_UP_ITERATION = 1;
            sound.volume(0.0);
            runningDiff = runningDiff + diff;
            var serverSongTime = data['song_time']+time_delta2/1000.0;
            console.log('[' + Date.now() + '] ' + ': NOT in sync (>' + MAX_SYNC_LAG.toString() + ' ms)')
            console.log('Browser:  ' + mySongTime.toString() + '\nServer: ' + serverSongTime.toString() + '\nDiff: ' + (diff*1000).toString() + '\nMean half-latency: ' + true_server_time_delta.toString() +  '\nMeasured half-latency: ' + time_delta2.toString() + '\nrunningDiff: ' + (runningDiff*1000).toString() + '\nSeeking to: ' + (serverSongTime+runningDiff).toString());
            $("div.info1").html('Muted <b>' + current_song_name + '</b> (out of sync)');
            if (diff<-1000000) {
              console.log('pausing')
              sound.pause()
              clearTimeout(secondTimeout3);
              clearTimeout(mainInterval);
              secondTimeout3 = setTimeout(function() {
                  console.log('playing');
                  sound.play();
                  mainInterval = setInterval(function(){
                    checkIfSkipped();
                  }, CHECK_UP_WAIT_TIME);
              }, Math.abs(runningDiff)*1000);
            } else {
                console.log(JSON.stringify(data));
                sound.seek(serverSongTime+runningDiff);
            }
          } else {
            console.log('[' + Date.now() + '] ' + ': in sync (|' + (diff*1000).toString() + '|<' + MAX_SYNC_LAG.toString() + ' ms)')
            $("div.info1").html('Playing <b>' + current_song_name + '</b>');
            CHECK_UP_ITERATION = parseInt(30.0/(CHECK_UP_WAIT_TIME/1000.0)); // every 30 seconds
            tryWait = 0;
            check_up_counter = 0;
            sound.volume(1.0);
          } 
        }
      }
    });

}




$(document).ready(function(){
$('a[type=controls]').click(function() {
   var skip = $(this).data('skip');
   console.log(skip);
    $("div.info1").text('Changing song');
    var posting = $.post('/nextsong', {
        'message': 'next song please',
        'skip': parseInt(skip)
    });

    // Put the results in a div
    posting.done(function(data) {
        sound.unload()
        location.reload(true);
        console.log('reloading page')
    });

});


makeRequests();

});




</script>

    <div class="container">
        <!-- columns should be the immediate child of a .row -->


        <div class="row">
        </div>


        <div class="row">
        <span style="display:table;">
<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px"
   width="138.32px" height="173.346px" viewBox="174.948 193.301 138.32 173.346"
   enable-background="new 174.948 193.301 138.32 173.346" xml:space="preserve">
<g transform="translate(-200.076,-237.528)">
  <path fill="#698000" stroke="#000000" stroke-width="2" d="M444.167,517.697c24.619-20.595,26.952,59.164,47.699,47.262
    c11.141-7.874,7.707-15.621,5.3-20.108c-3.354-6.256-6.615-8.804-1.917-8.588c3.574,0.163,20.693,14.632,5.898,29.888
    c-20.579,21.222-29.928,3.441-52.319,18.906C429.793,598.204,444.392,518.02,444.167,517.697L444.167,517.697z"/>
  <path fill-opacity="0.3137" d="M507.802,553.783c-1.141,2.051-2.518,5.787-4.598,7.932c-20.579,21.222-33.524,7.449-50.341,19.013
    c-7.743,5.326-9.914,5.458-9.795,3.701c-2.6,4.378-3.989,8.671,8.026,0.371c22.393-15.465,31.731,2.45,52.099-20.857
    C506.688,560.339,506.776,557.054,507.802,553.783z M493.723,536.328c-1.363,1.366,1.247,4.324,2.474,6.284
    c1.603,3.954,5.653,7.17-0.03,18.109C505.129,554.797,502.171,541.616,493.723,536.328L493.723,536.328z"/>
</g>
<path fill="#698000" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M283.128,310.215
  c0.006,25.82-20.925,46.756-46.744,46.755c-25.821-0.001-46.749-20.938-46.74-46.76c-0.007-25.82,20.924-46.754,46.744-46.753
  C262.209,263.458,283.137,284.396,283.128,310.215z"/>
<path fill-opacity="0.3137" d="M280.304,294.18c0.032,0.717-0.003,1.45,0.053,2.174c2.83,36.476-20.939,54.84-46.742,54.839
  c-20.168-0.001-39.93-18.749-43.925-38.814c1.131,24.799,21.621,44.581,46.698,44.582c25.801,0.002,46.74-20.936,46.742-46.737
  C283.13,304.589,282.133,299.183,280.304,294.18z"/>
<path fill="#698000" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M207.974,353.202
  c-0.003,7.402,5.998,8.346,13.399,8.347c7.401,0,13.401-0.943,13.399-8.346c0.003-7.401-6.248-11.488-13.398-13.404
  C210.252,336.818,207.972,345.8,207.974,353.202z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M231.804,357.378
  c0.537,2-1.061,2.701-3.115,3.251c-2.056,0.551-3.789,0.742-4.324-1.258c-0.536-2,0.695-4.067,2.75-4.618
  S231.269,355.377,231.804,357.378z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M218.28,359.544
  c-0.537,2-2.271,1.808-4.324,1.258c-2.055-0.55-3.653-1.252-3.116-3.252c0.536-2,2.636-3.175,4.691-2.624
  C217.584,355.476,218.816,357.544,218.28,359.544z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M224.898,359.214
  c0.001,2.071-1.724,2.335-3.851,2.335s-3.852-0.264-3.851-2.335c0.001-2.07,1.725-3.749,3.851-3.749
  C223.175,355.465,224.899,357.144,224.898,359.214L224.898,359.214z"/>
<path fill-opacity="0.3137" d="M225.616,341.265c1.832,2.403,2.958,5.403,2.957,9.108c0.002,8.754-6,9.872-13.402,9.871
  c-1.146,0-2.254-0.034-3.314-0.111c2.428,1.179,5.792,1.418,9.515,1.418c7.402,0,13.404-0.95,13.402-8.352
  C234.776,347.316,230.825,343.534,225.616,341.265z"/>
<path fill="#698000" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M264.8,353.205
  c0.001,7.402-5.998,8.345-13.4,8.344c-7.401,0-13.401-0.943-13.398-8.346c-0.003-7.401,6.249-11.488,13.398-13.402
  C262.523,336.821,264.802,345.804,264.8,353.205z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M262.156,357.38
  c0.537,2-1.062,2.7-3.116,3.25c-2.055,0.551-3.788,0.742-4.324-1.258s0.695-4.066,2.749-4.618
  C259.52,354.204,261.62,355.379,262.156,357.38L262.156,357.38z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M248.63,359.545
  c-0.534,2-2.268,1.809-4.323,1.258c-2.055-0.55-3.651-1.251-3.116-3.251c0.537-2,2.637-3.175,4.691-2.625
  C247.936,355.478,249.166,357.545,248.63,359.545z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M255.251,359.216
  c0,2.07-1.725,2.334-3.852,2.334s-3.852-0.264-3.851-2.334s1.724-3.75,3.851-3.75S255.251,357.146,255.251,359.216L255.251,359.216z
  "/>
<path fill-opacity="0.3137" d="M257.445,339.413c3.447,2.661,4.355,8.053,4.355,12.81c0.001,7.66-6.205,8.639-13.861,8.639
  c-1.662-0.001-3.246-0.054-4.723-0.204c2.264,0.725,5.099,0.888,8.18,0.888c7.402,0,13.404-0.94,13.402-8.342
  C264.8,347.404,263.4,340.641,257.445,339.413z"/>
<path fill="#698000" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M226.87,300.454
  c2.447,4.239,0.995,9.661-3.244,12.108c-4.24,2.448-9.662,0.996-12.11-3.244c-2.447-4.239-0.995-9.661,3.245-12.108
  C219,294.762,224.422,296.215,226.87,300.454z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M221.984,307.008
  c1.661,0.959,2.324,2.919,1.481,4.379c-0.842,1.46-2.871,1.864-4.532,0.905c-1.66-0.958-2.323-2.919-1.481-4.378
  C218.294,306.455,220.323,306.05,221.984,307.008z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M225.81,300.382
  c1.661,0.958,2.324,2.918,1.482,4.377c-0.843,1.46-2.872,1.864-4.533,0.906c-1.66-0.959-2.323-2.919-1.481-4.379
  C222.12,299.827,224.15,299.422,225.81,300.382z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M225.536,304.634
  c1.661,0.958,2.324,2.918,1.482,4.378c-0.844,1.459-2.872,1.864-4.533,0.905c-1.66-0.958-2.323-2.918-1.481-4.378
  C221.847,304.08,223.876,303.675,225.536,304.634z"/>
<path fill="#698000" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M261.257,309.318
  c-2.448,4.24-7.868,5.693-12.109,3.245c-4.239-2.448-5.692-7.869-3.244-12.109c2.448-4.239,7.869-5.692,12.109-3.244
  C262.252,299.657,263.705,305.079,261.257,309.318L261.257,309.318z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M251.082,307.177
  c-1.66,0.959-2.323,2.919-1.481,4.379c0.843,1.459,2.872,1.864,4.532,0.906c1.661-0.96,2.324-2.919,1.481-4.379
  C254.773,306.623,252.743,306.219,251.082,307.177z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M247.256,300.55
  c-1.661,0.959-2.324,2.919-1.481,4.378c0.842,1.459,2.871,1.865,4.531,0.906c1.661-0.959,2.324-2.919,1.482-4.378
  C250.947,299.996,248.917,299.591,247.256,300.55z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M247.53,304.802
  c-1.661,0.96-2.324,2.919-1.481,4.379c0.842,1.46,2.871,1.864,4.531,0.905c1.661-0.958,2.324-2.919,1.482-4.378
  S249.191,303.843,247.53,304.802L247.53,304.802z"/>
<path fill-opacity="0.3137" d="M259.868,243.24c0.55,3.002,1.137,6.186,1.137,9.491c0,17.773-17.295,32.2-25.677,32.2
  c-6.496,0-19.168-9.346-23.862-22.265c2.389,15.834,17.937,28.081,25.372,28.081c8.382,0,25.678-14.428,25.678-32.201
  C262.516,253.007,261.445,247.794,259.868,243.24z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linejoin="round" d="M243.599,275.85
  c-1.658,3.673,8.04,15.601,9.67,13.62c1.631-1.981,1.828-13.922-0.922-16.118C249.759,271.285,245.341,271.985,243.599,275.85z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linejoin="round" d="M229.175,275.85
  c1.657,3.673-8.041,15.601-9.67,13.62c-1.632-1.981-1.828-13.922,0.922-16.118C223.014,271.285,227.431,271.985,229.175,275.85z"/>
<path fill-opacity="0.3137" d="M229.173,275.845c1.657,3.673-8.044,15.61-9.676,13.629c-0.156-0.189-0.292-0.473-0.422-0.829
  c2.219,0.686,10.023-11.254,8.107-14.619c-0.263-0.46-0.569-0.87-0.896-1.226C227.47,273.378,228.511,274.381,229.173,275.845
  L229.173,275.845z"/>
<path fill="#698000" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M259.805,239.478
  c2.484-1.612,4.128-4.409,4.126-7.59c0.002-4.996-4.048-9.047-9.043-9.047c-2.616,0-4.971,1.112-6.622,2.888
  c-3.421-1.825-7.358-2.849-11.752-2.849c-4.537,0-8.548,1.025-11.997,2.857c-1.651-1.78-4.01-2.896-6.631-2.896
  c-4.996,0-9.044,4.051-9.044,9.047c0,3.194,1.657,6.001,4.158,7.61c-1.993,4.457-3.649,10.024-3.198,15.583
  c1.437,17.688,18.318,32.2,26.712,32.2c8.72,0,27.026-14.431,26.711-32.2C263.117,248.995,261.576,243.642,259.805,239.478z"/>
<path fill-opacity="0.3137" d="M260.476,239.768c0.573,3.003,1.184,6.186,1.184,9.492c0,17.773-17.992,32.199-26.712,32.199
  c-6.758,0-19.94-9.346-24.824-22.265c2.484,15.834,18.659,28.081,26.395,28.081c8.72,0,26.712-14.427,26.712-32.201
  C263.23,249.535,262.116,244.321,260.476,239.768z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M224.687,232.35
  c0.002,3.658-2.966,6.627-6.624,6.627c-3.661,0-6.626-2.968-6.626-6.627c0-3.66,2.965-6.627,6.626-6.627
  C221.721,225.723,224.689,228.689,224.687,232.35z"/>
<path d="M221.927,233.519c0.001,1.221-0.988,2.209-2.209,2.209c-1.219,0-2.209-0.988-2.208-2.209c0.001-1.22,0.988-2.21,2.208-2.21
  C220.939,231.309,221.928,232.299,221.927,233.519z"/>
<path fill="#FFFFFF" stroke="#000000" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" d="M248.086,232.35
  c-0.001,3.658,2.966,6.627,6.625,6.627s6.626-2.968,6.625-6.627c0.001-3.66-2.965-6.627-6.625-6.627
  C251.052,225.723,248.085,228.689,248.086,232.35z"/>
<path d="M250.847,233.519c0,1.221,0.989,2.209,2.208,2.209c1.221,0,2.21-0.988,2.208-2.209c0-1.22-0.987-2.21-2.208-2.21
  C251.835,231.309,250.847,232.299,250.847,233.519z"/>
<path d="M241.729,277.748c-0.748,1.297-2.121,1.905-3.065,1.358c-0.946-0.545-1.105-2.039-0.356-3.335
  c0.748-1.296,2.121-1.905,3.066-1.36C242.319,274.958,242.478,276.45,241.729,277.748z"/>
<path d="M231.044,277.748c0.749,1.297,2.122,1.905,3.067,1.358c0.945-0.545,1.104-2.039,0.355-3.335s-2.121-1.905-3.066-1.36
  C230.455,274.958,230.296,276.45,231.044,277.748z"/>
<g>
  <path fill="#00B0F0" stroke="#000000" d="M200.63,232.191h-4.069v-5.34h4.069V232.191z"/>
  <path fill="#00B0F0" stroke="#000000" d="M272.712,232.191h4.072v-5.34h-4.072V232.191z"/>
  <path fill="#0077BD" stroke="#000000" d="M208.639,226.427c0,0,5.74,7.176-0.718,23.68c0,0-16.503-0.719-15.068-8.611
    C192.853,241.495,192.857,226.427,208.639,226.427"/>
  <path fill="#EDBA54" stroke="#000000" d="M203.651,227.668c0,0,2.701,8.109-0.678,17.345c0,0,1.048-10.15-2.184-16.07
    C200.789,228.942,201.754,228.188,203.651,227.668"/>
  <path fill="#BEC0C2" stroke="#000000" d="M193.678,242.342c0,0-1.312-0.842,0.134-3.397c0,0-0.503-1.784-1.276-0.169
    c0,0-0.706,1.245-0.237,3.6C192.298,242.375,193.14,243.585,193.678,242.342"/>
  <path fill="#737577" stroke="#000000" d="M207.593,226.407c0,0,0.881-2.892,2.732-0.079c0,0,5.384,9.404-1.527,24.188
    c0,0-1.365,1.364-2.01-0.404C206.788,250.112,213.86,236.371,207.593,226.407"/>
  <path fill="#0077BD" stroke="#000000" d="M264.705,226.427c0,0-5.74,7.176,0.716,23.68c0,0,16.504-0.719,15.069-8.611
    C280.49,241.495,280.49,226.427,264.705,226.427"/>
  <path fill="#EDBA54" stroke="#000000" d="M269.694,227.668c0,0-2.702,8.109,0.676,17.345c0,0-1.045-10.15,2.184-16.07
    C272.553,228.942,271.59,228.188,269.694,227.668"/>
  <path fill="#BEC0C2" stroke="#000000" d="M279.666,242.342c0,0,1.313-0.842-0.136-3.397c0,0,0.505-1.784,1.278-0.169
    c0,0,0.706,1.245,0.237,3.6C281.046,242.375,280.205,243.585,279.666,242.342"/>
  <path fill="#737577" stroke="#000000" d="M265.751,226.407c0,0-0.883-2.892-2.73-0.079c0,0-5.386,9.404,1.524,24.188
    c0,0,1.366,1.364,2.01-0.404C266.555,250.112,259.483,236.371,265.751,226.407"/>
  <path fill="#0077BD" stroke="#000000" d="M199.131,226.786c0,0,33.365-59.736,76.687,1.347
    C275.818,228.133,238.417,158.977,199.131,226.786"/>
</g>
</svg>


            <span style="vertical-align: middle; display: table-cell;">
            <h1  style="position:relative;bottom:0"><i>musicsaur</i><br><small style="font-size: 50%%;">&nbsp;version 1.2</small></h1>
        </span>
    </div>


        <div class="row">
            <div class="seven columns">
                <a class="button" data-skip="-3" type="controls">Previous</a> <a class="button" data-skip="-2" type="controls">Replay</a> <a class="button" data-skip="-1" type="controls">Next</a>
            </div>


            <div class="five columns">
                <div class="info1">
                    {{ data['message'] }}
                </div>
            </div>
        </div>


        <div class="row">
            <div class="two columns">
            </div>


            <div class="ten columns">
                <div class="info2">
                    {{ data['playlist_html'] | safe }}
                </div>
            </div>
        </div>
    </div>
</body>
</html>

`
