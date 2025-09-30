# Module viam-utils

This module is used to be able to access methods that are only available in go through DoCommands.

## Model viam:viam-utils:arm

Go methods available in the arm RDK

### Configuration

The following attribute template can be used to configure this model:

```json
{
"arm": <string>,
}
```

#### Attributes

The following attributes are available for this model:

| Name  | Type   | Inclusion | Description                |
|------ |--------|-----------|----------------------------|
| `arm` | string | Required  | Name of the configured arm |

#### Example Configuration

```json
{
  "arm": "ur20"
}
```

### DoCommand

#### transform DoCommand

`transform` will transform the joint positions into a pose containing the translation and orientation vector.

```json
{
  "transform": {
    "joint_position": [
      0.1,
      0.2,
      0.3,
      0.4,
      0.5,
      0.6
    ],
  }
}
```
