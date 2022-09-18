<!DOCTYPE html>
<html>
  <head>
    <style>
      table {
        font-family: Arial, Helvetica, sans-serif;
        background-color: #EEEEEE;
        border-collapse: collapse;
        width: 100%;
      }
      table td, table th {
        border: 1px solid #ddd;
        padding: 3px 3px;
      }
      table th {
        font-size: 15px;
        font-weight: bold;
        padding-top: 12px;
        padding-bottom: 12px;
        text-align: left;
        background-color: #1C6EA4;
        color: white;
      }
    </style>
  </head>
  {{- $userPermissionMap := . }}
  {{- range $user, $objmap := $userPermissionMap }}
  <h1 style="font-size: 2em">{{ $user }}</h1>
  {{- range $object, $descriptor := $objmap.ObjectDesc }}
  <table>
    <tr>
      <th>{{ $object | Title }}/Permissions</th>
      {{- range $_, $perm := index $objmap.Permissions $object }}
      <th>{{ $perm }}</th>
      {{- end }}
    </tr>
    {{- range $item, $_ := $descriptor }}
      <tr>
        <td>{{ $item }}</td>
      {{- range $_, $status := index $objmap.Permissions $object }}
        {{- $flag := index $descriptor $item $status }}
        {{- if $flag }}
        <td style="background-color: green"> {{ $flag }}</td>
        {{- else }}
        <td style="background-color: red"> {{ $flag }}</td>
        {{- end }}
      {{- end }}
      </tr>
    {{- end }}
  </table>
  <table>
    <tr>
      <th>{{ $object | Title }}/Failed to check</th>
    </tr>
    {{- range $skipped_item_name, $skipped_item_type := $objmap.FailedToCheck }}
      {{- if eq $skipped_item_type $object }}
      <tr>
        <td style="background-color: yellow">{{ $skipped_item_name }}</td>
      </tr>
      {{- end }}
    {{- end }}
    <br>
  </table>
  <br>
  {{- end }}
{{- end }}
</html>
