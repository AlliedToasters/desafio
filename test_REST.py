import __future__
import os
import json
import os

def call(
        verb,
        errors,
        *args,
        POST=None,
        id=None,
        url='http://localhost:8080',
        uri_base='/item',
        query_params={},
        I = False
    ):
    """Executes REST call and returns output
    verb = GET, POST, or DELETE (string)
    args = any arguments for uri (strings)
    id = integer id to append to call
    url = base url for call
    uri_base = uri signiture
    kwargs = arguments to POST
    """
    call = "curl -s "
    if I:
        call += "-I "
    call += '-X "{}"'.format(verb)
    if POST != None:
        call += " -d '"
        call += json.dumps(POST) + "'"

    call += ' {}{}'.format(url, uri_base)
    if id != None:
        call += '/' + str(id)
    if len(query_params) > 0:
        call += '\?'
        for param in query_params:
            call += param + "="
            call += query_params[param] + "&"
    print('\nCall: \n', call, '\n')
    response = os.popen(call).read()
    print('\nResponse: \n', response, '\n')
    if not I:
        try:
            result = json.loads(response)
        except json.decoder.JSONDecodeError:
            print('Invalid JSON returned...')
            errors += 1
            return None, errors
    elif I:
        result = response
    return result, errors

def exit_test(errors):
     """Prints a report before calling quit()"""
     print('\n Test script run with {} errors.'.format(errors))
     quit()
     return


if __name__ in '__main__':
    errors = 0

    names = [
        'Mary',
        'Jose',
        'Nimbus A.'
    ]

    descriptions = [
        'a woman',
        'hombre, cliente',
        'a cat'
    ]

    ids = []

    response, errors = call('GET', errors, id=-1)
    if response == None:
        print('Script exited early with {} errors.'.format(errors))
        print('Did not recieve JSON response from server.')
        quit()
    try:
        assert 'error' in response
    except AssertionError:
        errors += 1
        print('Expected error response. Got: {}'.format(response))

    for n in range(3):
        post = {'name':names[n], 'description':descriptions[n]}
        response, errors = call('POST', errors, POST=post)
        ids.append(response['id'])

    for i, id in enumerate(ids):
        response, errors = call('GET', errors, id=id)
        try:
            assert response['name'] == names[i]
        except AssertionError:
            errors += 1
            print('Expected name {} at id {}. Got: {}'.format(
                names[i],
                id,
                response['name']
            ))
    try:
        assert response['description'] == descriptions[i]
    except AssertionError:
        errors += 1
        print('Expected description {} at id {}. Got: {}'.format(
            descriptions[i],
            id,
            response['description']
        ))

    table, errors = call('GET', errors)
    results = dict()
    results['id'] = []
    results['name'] = []
    results['description'] = []
    for item in table:
        results['id'].append(item['id'])
        results['name'].append(item['name'])
        results['description'].append(item['description'])

    idxs = []
    for i, id in enumerate(ids):
        idx = results['id'].index(str(id))
        idxs.append(idx)
        try:
            assert results['name'][idx] == names[i]
        except:
            errors += 1
            print('Expected name {} at id {}. Got: {}'.format(
                names[i],
                id,
                results['name'][idx]
            ))
        try:
            assert results['description'][idx] == descriptions[i]
        except:
            errors += 1
            print('Expected description {} at id {}. Got: {}'.format(
                descriptions[i],
                id,
                results['description'][idx]
            ))

    #Test DELETE on Jose (position 1):
    response, errors = call('DELETE', errors, id=ids[1])
    try:
        assert response['name'] == 'Jose'
    except:
        errors += 1
        print('Unexpected response from DELETE call: ', response)

    table, errors = call('GET', errors)
    results = dict()
    results['id'] = []
    for item in table:
        results['id'].append(item['id'])
    try:
        assert ids[1] not in results['id']
    except:
        errors += 1
        print('Error: Found item after deletion.')

    #Test DELETE call again, expect error response:
    response, errors = call('DELETE', errors, id=ids[1])
    try:
        assert 'error' in response
    except:
        errors += 1
        print('Expected error response from DELETE call, got: ', response)

    #Test the Google Drive API DB interface
    response, errors = call('GET', errors, id=1, uri_base='/file')
    try:
        auth_url = response['description']
        auth_code = input("Go to URL and paste code here: \n{}\n".format(auth_url))
        post = {'auth_code':auth_code}
        response, errors = call('POST', errors, POST=post, uri_base='/auth')
        try:
            if "success" not in response:
                errors += 1
                print("Authentication failed. Exiting...")
                exit_test(errors)
            else:
                print("Authentication success. Waiting for database to populate...")
                os.system("sleep 5")
        except TypeError:
            errors += 1
            print("Authentication failed. Exiting...")
            exit_test(errors)
    except KeyError:
        try:
            response["id"]
        except KeyError:
            errors += 1
            print("unexpected response: ", response)



    response, errors = call('GET', errors, uri_base='/file')

    response, errors = call('GET', errors, id=1, uri_base='/file')

    query = {"word":"dev"}
    response, errors = call(
        'GET',
        errors,
        id=1,
        uri_base='/search-in-drive',
        query_params=query,
        I=True
    )
    code = response.split(" ")[1]
    if code != "200":
        errors += 1
        print("Error: Expected code response 200. Got: ", code)

    query = {"word":"honk"}
    response, errors = call(
        'GET',
        errors,
        id=1,
        uri_base='/search-in-drive',
        query_params=query,
        I=True
    )
    code = response.split(" ")[1]
    if code != "404":
        errors += 1
        print("Error: Expected code response 404. Got: ", code)

    exit_test(errors)
