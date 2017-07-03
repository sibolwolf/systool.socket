package audioctrller

import (
	"errors"
	"io"
	"time"

	"smartconn.cc/liugen/audio"
	sysLocker "smartconn.cc/sibolwolf/syssleepwake"
	"smartconn.cc/sibolwolf/syssleepwake/sleephandle"
	"smartconn.cc/sibolwolf/syssleepwake/wakehandle"
	"smartconn.cc/tosone/logstash"
)

var isBreak = true // 当前正在播放是 break 或者是 play

// 待播放文件的信息
type audioFile struct {
	file string
	time int64
}

var list []audioFile

var cacheFile audioFile // 当前正在的 play 或者被播放到一半的 play

var timer int64

func init() {
	list = []audioFile{}
	audio.Startup()

	// Register audio close before sleep
	sleephandle.StopAudioEventRegister(TearDown)

	// Register audio open after wake
	wakehandle.StartAudioEventRegister(func() {
		err := Initialize()
		if err != nil {
			logstash.Error(err.Error())
		}
	})

}

func countTimer() {
	timer = 0
	for audio.IsBGMPlaying() {
		timer++
		<-time.After(time.Second)
	}
}

var appName = "audioCtrl"

func lock() {
	sysLocker.UpdateLockStatus(appName, sysLocker.Lock)
}

func unlock() {
	sysLocker.UpdateLockStatus(appName, sysLocker.Unlock)
}

// Play 异步播放，不会卡住主线程
func Play(file string) error {
	var (
		err     error
		channel chan bool
	)
	if audio.IsBGMPlaying() {
		if !isBreak {
			list = append(list, audioFile{file: file})
			return err
		}
		audio.StopBGM()
	}

	isBreak = false
	go countTimer()
	cacheFile = audioFile{file: file}
	go func() {
		logstash.WithFields(logstash.Fields{"audio": file}).Info("Playing audio.")
		channel, err = audio.PlayBGM(file)
		lock()
		if err != nil {
			return
		}
		<-channel
		cacheFile = audioFile{}
		Resume()
		unlock()
	}()
	<-time.After(time.Millisecond * 200)
	return err
}

// PlaySync 同步播放，主线程将会被卡住
func PlaySync(file string) error {
	var (
		err     error
		channel chan bool
	)

	if audio.IsBGMPlaying() {
		if !isBreak {
			return errors.New("Audio is busy")
		}
		audio.StopBGM()
	}

	isBreak = false
	go countTimer()
	cacheFile = audioFile{file: file}
	channel, err = audio.PlayBGM(file)
	lock()
	if err != nil {
		return err
	}
	<-channel
	cacheFile = audioFile{}
	Resume()
	unlock()
	return err
}

// PlayAt 异步播放，不会卡住主线程
func PlayAt(file string, startAt int64) error {
	var (
		err     error
		channel chan bool
	)

	if audio.IsBGMPlaying() {
		if !isBreak {
			list = append(list, audioFile{file: file})
			return err
		}
		audio.StopBGM()
	}

	isBreak = false
	go countTimer()
	cacheFile = audioFile{file: file}
	go func() {
		channel, err = audio.PlayBGM(file)
		if err != nil {
			return
		}
		lock()
		<-channel
		cacheFile = audioFile{}
		Resume()
		unlock()
	}()
	<-time.After(time.Millisecond * 200)
	return nil
}

// PlayAtSync 同步播放，主线程将会被卡住
func PlayAtSync(file string, startAt int64) error {
	var (
		err     error
		channel chan bool
	)

	if audio.IsBGMPlaying() {
		if !isBreak {
			return errors.New("Audio is busy")
		}
		audio.StopBGM()
	}

	isBreak = false
	go countTimer()
	cacheFile = audioFile{file: file}
	channel, err = audio.PlayBGM(file, startAt)
	if err != nil {
		return err
	}
	lock()
	<-channel
	cacheFile = audioFile{}
	Resume()
	unlock()
	return err
}

// Break break 播放
func Break(file string) error {
	var (
		err     error
		channel chan bool
	)

	isBreak = true
	Pause()
	go func() {
		channel, err = audio.PlayBGM(file)
		if err != nil {
			Resume()
		}
		lock()
		<-channel
		Resume()
		unlock()
	}()
	return nil
}

// BreakSync break 播放
func BreakSync(file string) error {
	var (
		err     error
		channel chan bool
	)

	isBreak = true
	channel, err = audio.PlayBGM(file)
	if err != nil {
		Resume()
		return err
	}
	lock()
	<-channel
	Resume()
	unlock()
	return err
}

// PlaySE SE 播放
func PlaySE(src interface{}) error {
	channel, err := audio.PlaySE(src)
	if err != nil {
		return err
	}
	lock()
	<-channel
	unlock()
	return err
}

// Pause 暂停正在的播放
func Pause() {
	if IsPlaying() {
		audio.StopBGM()
		if !isBreak {
			cacheFile = audioFile{
				file: cacheFile.file,
				time: cacheFile.time,
			}
		}
	}
	unlock()
}

// Clear 清理播放队列
func Clear() {
	list = []audioFile{}
	audio.StopBGM()
	audio.StopSE()
	audio.StopRecord()
	cacheFile = audioFile{}
	unlock()
}

// Resume 恢复之前的播放
func Resume() {
	if cacheFile.file != "" {
		PlayAt(cacheFile.file, cacheFile.time)
	} else if audioInfo := queue(); audioInfo.file != "" {
		Play(audioInfo.file)
	}
}

// IsPlaying 是否正在播放
func IsPlaying() bool {
	return audio.IsBGMPlaying()
}

// queue 获取队列中的第一个值
func queue() audioFile {
	if len(list) == 0 {
		return audioFile{}
	}
	audioInfo := list[0]
	if len(list) == 1 {
		list = []audioFile{}
	} else {
		list = list[1:]
	}
	return audioInfo
}

// Record 开始录音
func Record() (io.Reader, error) {
	lock()
	return audio.Record()
}

// StopRecord 停止录音
func StopRecord() {
	audio.StopRecord()
	unlock()
}

// IsRecording 是否在录音
func IsRecording() bool {
	return audio.IsRecording()
}

// TearDown 卸载声卡
func TearDown() {
	audio.Teardown()
	unlock()
}

// Initialize 初始化
func Initialize() error {
	return audio.Startup()
}

// Convert 转换音频格式
func Convert(src, dest string) error {
	return audio.Convert(src, dest)
}
