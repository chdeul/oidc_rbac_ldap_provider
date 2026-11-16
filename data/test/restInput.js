//  var input = json bytes 

function jsExec( json ) {

/*  input is 
 {
  "users": [
    {
      "id": 1,
      "name": "administrator",
      "real_name": "Super Admin",
      "email": "root@localhost",
      "access_level": {
        "id": 90,
        "name": "administrator",
        "label": "administrateur"
      }
    }, ....
  ]
} 
computed to 
  {
      "Uid": 1,
      "Identifiant": "administrator",
      "NomComplet": "Super Admin",
      "AdresseMail": "root@localhost",
      "Roles": "90:administrateur"
  }
*/
  obj = JSON.parse(json)

  if ( obj && obj.users && obj.users.length > 0 ) {
    const u0 = obj.users[0]
    user = { 
        "Uid": u0.id,
        "Identifiant": u0.name,
        "NomComplet": u0.real_name,
        "AdresseMail": u0.email,
        "Roles": "" + u0.access_level.id + ":" + u0.access_level.label, 
    }
    // return json bytes
    return JSON.stringify(user)
  } 
  else { 
    return "error"
  }
}
 
// Do the job
jsExec(input)
