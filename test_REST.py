import __future__
import os
import json

def call(
        verb,
        errors,
        *args,
        POST=None,
        id=None,
        url='http://localhost:8080',
        uri_base='/item',
    ):
    """Executes REST call and returns output
    verb = GET, POST, or DELETE (string)
    args = any arguments for uri (strings)
    id = integer id to append to call
    url = base url for call
    uri_base = uri signiture
    kwargs = arguments to POST
    """
    call = 'curl -s -"X" {}'.format(verb)
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
    try:
        result = json.loads(response)
    except json.decoder.JSONDecodeError:
        print('Invalid JSON returned...')
        errors += 1
        return None, errors
    return result, errors

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
    titulos = [
        'documento1.docx',
        'documento2.txt',
        'image0.png'
    ]

    descripciones = [
        'a formatted text document',
        'a raw text file',
        'an image file'
    ]

    ids = []

    response, errors = call('GET', errors, id=-1, uri_base='/file')
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
        post = {'titulo':titulos[n], 'descripcion':descripciones[n]}
        response, errors = call('POST', errors, POST=post, uri_base='/file')
        ids.append(response['id'])

    for i, id in enumerate(ids):
        response, errors = call('GET', errors, id=id, uri_base='/file')
        try:
            assert response['titulo'] == titulos[i]
        except AssertionError:
            errors += 1
            print('Expected titulo {} at id {}. Got: {}'.format(
                titulo[i],
                id,
                response['titulo']
            ))
    try:
        assert response['descripcion'] == descripciones[i]
    except AssertionError:
        errors += 1
        print('Expected descripcion {} at id {}. Got: {}'.format(
            descripciones[i],
            id,
            response['descripcion']
        ))

    table, errors = call('GET', errors, uri_base='/file')
    results = dict()
    results['id'] = []
    results['titulo'] = []
    results['descripcion'] = []
    for documento in table:
        results['id'].append(documento['id'])
        results['titulo'].append(documento['titulo'])
        results['descripcion'].append(documento['descripcion'])

    idxs = []
    for i, id in enumerate(ids):
        idx = results['id'].index(str(id))
        idxs.append(idx)
        try:
            assert results['titulo'][idx] == titulos[i]
        except:
            errors += 1
            print('Expected name {} at id {}. Got: {}'.format(
                titulos[i],
                id,
                results['titulos'][idx]
            ))
        try:
            assert results['descripcion'][idx] == descripciones[i]
        except:
            errors += 1
            print('Expected descripcion {} at id {}. Got: {}'.format(
                descripciones[i],
                id,
                results['descripcion'][idx]
            ))




    print('\n Test script run with {} errors.'.format(errors))
