import __future__
import os
import json

def call(verb, *args, POST=None, id=None, url='http://localhost:8080', uri_base='/item'):
    """Executes REST call and returns output
    verb = GET, POST, or DELETE (string)
    args = any arguments for uri (strings)
    id = integer id to append to call
    url = base url for call
    uri_base = uri signiture
    kwargs = arguments to POST
    """
    call = 'curl -"X" {}'.format(verb)
    if POST != None:
        call += " -d '"
        call += json.dumps(POST) + "'"
        print(call)

    call += ' {}{}'.format(url, uri_base)
    print('\nCall: \n', call, '\n')
    if id != None:
        call += '/' + str(id)
    response = os.popen(call).read()
    print('\nResponse: \n', response, '\n')
    return json.loads(response)

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

    response1 = call('GET', id=1)
    try:
        assert 'error' in response1
    except AssertionError:
        errors += 1
        print('Expected error response. Got: {}'.format(response))

    for n in range(3):
        post = {'name':names[n], 'description':descriptions[n]}
        response = call('POST', POST=post)
        ids.append(response['id'])

    for i, id in enumerate(ids):
        response = call('GET', id=id)
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

    table = call('GET')
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
    response = call('DELETE', id=ids[1])
    try:
        assert response['name'] == 'Jose'
    except:
        errors += 1
        print('Unexpected response from DELETE call: ', response)

    table = call('GET')
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
    response = call('DELETE', id=ids[1])
    try:
        assert 'error' in response
    except:
        errors += 1
        print('Expected error response from DELETE call, got: ', response)

    print('\n Test script run with {} errors.'.format(errors))
