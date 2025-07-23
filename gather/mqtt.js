const mqtt = require('mqtt');
 
const client = mqtt.connect('ws://127.0.0.1:1882');
const _topic = 'fastsearch/';
client.on('connect', function() {
    console.log('Connected to MQTT Broker');
    setInterval(() => {
        publish();
    }, 1);
});
// client.subscribe(_topic+'index/back'); 
client.subscribe(_topic+'index/err'); 
client.on('message', function(topic, message) {
    // message is Buffer
    console.log("收到消息",message.toString(),"来自",topic);
});
 
client.on('error', function(err) {
    console.log("ERROR",err);
});

client.on('reconnect', function () {
    console.log('正在尝试重新连接...');
  });
   
  client.on('close', function () {
    console.log('连接已关闭');
    client.reconnect();
});
   
  //发布消息
 let index=1
  function publish(){
    const message={
        id:'sn-'+index,
        text:'text'+index,
        title:'title'+index,
        cut_document:true,
        has_key:true,
        keys:{
            id:'sn-'+index,
            name:'test',
            site:'test.com',
        },
        document:{
            name:'test'+index,
            age:18,
        },
        tags:['test','fastsearch'],
        
    }
    let msg=JSON.stringify(message)
    client.publish(_topic+"index",msg,{qos : 2})
    index++
  }