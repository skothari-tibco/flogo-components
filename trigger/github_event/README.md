---
title: Github Trigger

---
# tibco-gihub_trigger
This trigger provides your flogo application the ability to listen for a Github Event


## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings": [
      {
        "name": "port",
        "type": "integer",
        "required": true
      },
      {
          "name":"url",
          "type":"string",
          "required":true
      }
    ],
    "output": [
        {
            "name":"data",
            "type":"any"
        }
        
    ],
    "handler":{
        "settings":[
         {
          "name": "get_all",
          "type": "boolean",
          "required" : true
        }
         
        ]
    }
}
```
## Settings
### Trigger:
| Setting     | Description    |
|:------------|:---------------|
| port | The port to listen on | 
| url | The url to listen on. This will be a webhook you add to you Repository |            
### Endpoint:
| Setting     | Description    |
|:------------|:---------------|
| get_all      | Flag to indicate what type of data to recieve. |         





