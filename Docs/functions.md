# Functions

TODO function environments

TODO terminology: function environemtn (fnenv), function, function reference, function identifier, resolver

TODO anatomy of a function reference

There are currently two function environments: **Fission** and **Internal**.

[Fission](https://github.com/fission/fission) is a complete Function-as-a-Service platform - which includes extensive
Kubernetes integration, autoscaling, and offers fine-grained resource management controls.

The **internal** function environment is an extremely light-weight function runtime built into the 
workflow engine itself.
It is useful for lightweight functions - typically control flow functions - where the network overhead of running 
them on an external FaaS platform would be overkill.
However, it has no sandboxing or autoscaling options; intensive functions might slow down the workflow engine.

Whether it is best to execute your functions on either of the function environments is based on the kind of function 
and the type of workload.

If a function happens to be declared on both platforms for now the workflow engine will always choose the internal 
function to use.
You can force the workflow engine to use one fnenv over the other by specifying it in the function reference.
For example, to explicitly use a Fission function called `sleep`, the function reference should: `fission://sleep`.   

In the future we want to improve the scheduler to allow it to make this decision to based on profiling and other data 
sources. 

## Fission

[Fission](https://github.com/fission/fission) is the Function-as-a-Service (FaaS) platform underlying the workflow 
engine.
All Fission functions are compatible with the workflow engine.
Moreover, workflows follow the same API - allowing you to call Fission functions without even needing to know whether 
it is a workflow or an actual Fission function.

### Specification

Calling a Fission function is no different from calling other (internal) functions.
Like other functions, it has a number of optional input parameters to augment the Fission function execution.

**input**        | required | types                | description
-----------------|----------|----------------------|---------------------------------
body/default     | no       | *                    | The body of the Fission function call.
headers          | no       | map[string]string    | The headers that need to be added to the request.
query            | no       | map[string]string    | The key-value values that need to be added to the URL.
method           | no       | string               | HTTP Method to use (Default: GET)
content-type     | no       | string               | Force a specific content-type for the request (default: application/octet-stream).

**output** (*) the body of the Fission function response.

The body is interpreted based on the content-type.
In case it is interpretable (e.g. `application/json`), the workflow engine will parse it, and make it available to 
the [expressions](./expressions.md) to manipulate/access.
In case it is not interpretable, the output will be considered a binary blob.
You can still pass it as a whole to other tasks, but you can not access or modify the contents with expressions. 

Note: currently you cannot access the metadata of the response. Though, You can access this metadata for the workflow 
invocation.

**Fission's Perspective**

The workflow engine invokes Fission functions with the same API as other types of event sources.
So conceptually there is no difference between Fission functions used or not used in workflows.

However, for the sake of debugability and traceability, Workflows adds a bit of Workflows-specific metadata to each 
function.
All of this metadata is optional and not required (or expected) for correct execution of the workflows. 

**header**       | description
-----------------|-------------
Fission-Function | TODO

### Notes
- The content-type is important if you want to utilize the full functionality of Workflows; ensure that the functions 
have the correct MIME/content type in their responses.
- There is no need to force content-types on requests in many cases. 
The workflow engine remembers the content-type of how it received the data, and will use it when using the data as 
input for another task (unless the content type is overriden).
- You can return a specification for a task or workflow (to implement dynamic tasks) by using the appropriate 
content-type: `application/vnd.fission.workflows.task` or `application/vnd.fission.workflows.workflow` using the 
protobuf encoding.

## Internal

TODO intro 

It consists out of a number of built-in functions, which aim to cover the common functionality needed in workflows.

### Built-in

The internal fnenv ships with a number of built-in functions. 
These functions are simply commonly used, utility functions. 
They do not have any additional or special API; you could implement these functions just as well in Fission.

#### compose

Property  | description
----------|--------
command   | `compose`
available | `^0.1.1`

**Description**

Compose

**Use Cases** ?

**Specification**
TODO input
TODO output

**Example**
TODO short partial example

TODO Link to runnable example

#### fail

Property  | description
----------|--------
command   | `fail`
available | `^0.3.0`
status    | partial (custom message not yet implemented)

**Description**

Fail is a function that always fails. This can be used to short-circuit workflows in
specific branches. Optionally you can provide a custom message to the failure.

**Specification**

**input**   | required | types  | description
------------|----------|--------|---------------------------------
default     | no       | string | custom message to show on error 

**output** None 

**Example**

```yaml
# ...
foo:
  run: fail
  inputs: "all has failed"
# ...
```

A complete example of this function can be found in the [failwhale](../examples/whales/failwhale.wf.yaml) example.

#### foreach

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

**Description**

**Specification**

**Example**


#### http

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

**Description**

**Specification**

**Example**

#### if

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

**Description**

**Specification**

**Example**


#### javascript

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

#### noop

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

#### repeat

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

#### sleep

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

#### switch

Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

#### while
 
Property  | description
----------|--------
command   | `while`
available | `^0.3.0`

### Extending the internal function environment
As the naming suggests, the internal function environment is not limited to this small set of built-in functions.
The internal function environment can be extended with functions just like Fission.

To add an additional function, you need to implement the `native.Function` Go interface, and add the new function to
the list of functions.

Currently, you will still need to recompile the engine if you want add, change or remove internal functions.
In the near future you will be able to add these functions (in Go or Javascript) to the workflow engine at runtime.
